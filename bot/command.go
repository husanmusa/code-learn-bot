package bot

import (
	"database/sql"
	"errors"
	"github.com/husanmusa/code-learn-bot/pkg/parser"
	"github.com/husanmusa/code-learn-bot/pkg/structs"
	"github.com/husanmusa/code-learn-bot/service/lesson"
	"github.com/husanmusa/code-learn-bot/service/user"
	tele "gopkg.in/telebot.v3"
	"log"
	"strconv"
	time "time"
)

const userCtxKey = "user"

func handleStart(ctx tele.Context) error {
	return ctx.Send(enterName)
}

func handlerGetName(userService user.Service) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		var (
			_    = ctx.Sender()
			text = ctx.Text()
			user = ctx.Get(userCtxKey).(*structs.User)
		)
		if user.InputName == "" {
			user.InputName = text
			if err := userService.Update(user.ID, user); err != nil {
				return err
			}
		}
		menu.Reply(
			menu.Row(btnChoose),
		)
		if !sendingTask {
			return ctx.Reply(textHello, menu)
		} else {
			return handlerForwardMessage(userService)(ctx)
		}
	}
}

func handlerInfo(ctx tele.Context) error {
	menu.Reply(
		menu.Row(btnLessons),
	)
	return ctx.Reply(texInfo, menu)
}

func handlerLessons(lessonService lesson.Service, userService user.Service) tele.HandlerFunc {
	return func(ctx tele.Context) error {

		user := ctx.Get(userCtxKey).(*structs.User)
		user, err := userService.ReadByChatID(user.ID)
		var id = user.DoingLesson

		lessons, err := lessonService.ReadByLessonID(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = ctx.Reply(strconv.FormatInt(id, 10) + "-raqamli dars topilmadi")
				if err != nil {
				}
				log.Println(err)
			} else {
				log.Println(err)
			}
		}

		messages, _ := parser.ParseLessonToMessage(lessons)
		for _, message := range messages {
			edit := Editable{strconv.FormatInt(int64(message.ID), 10), message.Chat.ID}
			err = ctx.Forward(edit)
			if err != nil {
				log.Println("Error in Forward", err)
			}
		}

		user.LessonTime = time.Now()
		err = userService.Update(user.ID, user)
		if err != nil {
			log.Println(err)
		}
		menu.Reply(
			menu.Row(btnTask),
		)

		return ctx.Send(parser.ParseTimeToMessage(user.LessonTime), menu)
	}
}

func handleGetLessons(lessonService lesson.Service) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		msg := ctx.Message()

		parsed, err := parser.ParseMessageToLesson(msg)

		if err != nil {
			log.Println(err)
		}
		err = lessonService.Store(&parsed)
		if err != nil {
			log.Println(err)
			return err
		}

		return ctx.Reply("Added new lesson")
	}
}

func handleTask(userService user.Service) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		user := ctx.Get(userCtxKey).(*structs.User)
		user, err := userService.ReadByChatID(user.ID)
		if err != nil {
			log.Println(err)
		}
		startedTime := time.Now().UTC().Add(time.Hour * 5).Sub(user.LessonTime)
		duration, err := time.ParseDuration(taskTime)
		if err != nil {
			log.Println(err)
		}
		if startedTime-duration < 0 {
			return ctx.Reply(parser.ParseTimeToMessage(user.LessonTime))
		} else {
			//
			//err = userService.Update(user.ID, user)
			//if err != nil {
			//	log.Println(err)
			//}
			sendingTask = true
			return ctx.Reply(textSendTask)
		}
	}
}

func handlerForwardMessage(userService user.Service) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		user := ctx.Get(userCtxKey).(*structs.User)
		user.DoingLesson++
		err := userService.Update(user.ID, user)
		if err != nil {
			log.Println(err)
		}
		sendingTask = false
		var rcp tele.Recipient = Recipient{}

		err = ctx.ForwardTo(rcp)
		if err != nil {
			return err
		}
		menu.Reply(
			menu.Row(btnNextLesson),
		)
		return ctx.Send(isNext, menu)
	}
}

func handlerNextLesson(lessonService lesson.Service, userService user.Service) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		return handlerLessons(lessonService, userService)(ctx)
	}
}

type Editable struct {
	messageID string
	chatID    int64
}

func (e Editable) MessageSig() (messageID string, chatID int64) {
	return e.messageID, e.chatID
}
