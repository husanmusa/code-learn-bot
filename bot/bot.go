package bot

import (
	"context"
	"errors"
	"fmt"
	"github.com/husanmusa/code-learn-bot/pkg/parser"
	"github.com/husanmusa/code-learn-bot/service/lesson"
	"github.com/husanmusa/code-learn-bot/service/user"
	"log"
	"time"

	tele "gopkg.in/telebot.v3"
)

var (
	// Universal markup builders.
	menu          = &tele.ReplyMarkup{ResizeKeyboard: true, RemoveKeyboard: true}
	started       bool
	taskTime      = "1m"
	sendingTask   bool
	sendFeedback  bool
	btnNextLesson = menu.Text("🧠 Начать следующую задачу!")
	btnChoose     = menu.Text("🤯 Бесплатно, но самостоятельно")
	btnCash       = menu.Text("Платно - 10$, с куратором")
	btnTask       = menu.Text("✅ Задача готова! Хочу загрузить результат!")
	btnLessons    = menu.Text("👌 Я все понял, начать изучение")
	btnCheckPaid  = menu.Text("☑️Check for paid")
	btnSummary    = menu.Text("оставить отзыв")
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
	bot.Handle(&btnCash, handlerCasher)
	bot.Handle(&btnSummary, func(ctx tele.Context) error {
		sendFeedback = true
		menu.Reply()
		return ctx.Send("Type your feedbacks", menu)
	})

	withLesson.Handle(&btnLessons, handlerLessons(bot, lessonService, userService))
	withLesson.Handle(&btnTask, handleTask(userService, lessonService))
	withLesson.Handle(tele.OnChannelPost, handleGetLessons(lessonService))
	withLesson.Handle(&btnNextLesson, handlerLessons(bot, lessonService, userService))

	//bot.Handle(tele.OnText, handlerForwardMessage(userService))
	withUser.Handle(&btnCheckPaid, handlerCheckPaid)
	withUser.Handle(tele.OnText, handlerGetName(userService))

	bot.OnError = func(err error, ctx tele.Context) {
		if errors.Is(err, parser.ErrorLessonWrong) {
			if e := ctx.Reply(fmt.Sprintf("%s", err)); e != nil {
				log.Println("Error Reply: ", e)
			}
		} else if errors.Is(err, parser.ErrorLessonThanNeed) {
			if e := ctx.Reply(fmt.Sprintf("%s", err)); e != nil {
				log.Println("Error Reply 2: ", e)
			}
		} else {
			if e := ctx.Reply("Что-то пошло не так. Попробуйте позже или обратитесь к @HusanMusa"); e != nil {
				log.Println(e)
			}
		}

		log.Println(err)
	}
	bot.Start()

	<-ctx.Done()
	bot.Stop()

	return nil
}

type Recipient struct {
	check bool
}

func (r Recipient) Recipient() string {
	if r.check {
		return "-571389424"
	}
	return "-1001371913344"
}

type Recip struct {
	id string
}

func (r Recip) Recipient() string {
	return r.id
}
