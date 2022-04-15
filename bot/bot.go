package bot

import (
	"context"
	"github.com/husanmusa/code-learn-bot/service/lesson"
	"github.com/husanmusa/code-learn-bot/service/user"
	"log"
	"time"

	tele "gopkg.in/telebot.v3"
)

var (
	// Universal markup builders.
	menu = &tele.ReplyMarkup{ResizeKeyboard: true, RemoveKeyboard: true}
	//selector = &tele.ReplyMarkup{}
	taskTime      = "1m"
	sendingTask   bool
	btnNextLesson = menu.Text("🧠 Начать следующую задачу!")
	btnChoose     = menu.Text("🤯 Бесплатно, но самостоятельно")
	btnTask       = menu.Text("✅ Задача готова! Хочу загрузить результат!")
	btnLessons    = menu.Text("👌 Я все понял, начать изучение")
	//
	//btnPrev = selector.Data("⬅", "prev", "Te")
	//btnNext = selector.Data("➡", "next", "st")
)

func Start(ctx context.Context, token string, userService user.Service, lessonService lesson.Service,
) error {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("connected bot %q\n", bot.Me.Username)

	var (
		authUser   = auth(userService)
		withUser   = bot.Group()
		withLesson = bot.Group()
	)
	withLesson.Use(authUser)
	withUser.Use(authUser)
	bot.Handle("/start", handleStart)
	bot.Handle(&btnChoose, handlerInfo)

	withLesson.Handle(&btnLessons, handlerLessons(lessonService, userService))
	withLesson.Handle(&btnTask, handleTask(userService))
	withLesson.Handle(tele.OnChannelPost, handleGetLessons(lessonService))
	withLesson.Handle(&btnNextLesson, handlerLessons(lessonService, userService))

	bot.Handle(tele.OnText, handlerForwardMessage(userService))

	withUser.Handle(tele.OnText, handlerGetName(userService))

	bot.OnError = func(err error, ctx tele.Context) {
		if e := ctx.Reply("Что-то пошло не так. Попробуйте позже или обратитесь к @HusanMusa"); e != nil {
			log.Println(e)
		}

		log.Println(err)
	}
	bot.Start()

	<-ctx.Done()
	bot.Stop()

	return nil
}

type Recipient struct {
}

func (Recipient) Recipient() string {
	return "-1001371913344"
}
