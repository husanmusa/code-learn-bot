package bot

import (
	"database/sql"
	"errors"
	"fmt"
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
	started = true
	return ctx.Send(enterName)
}

func handlerGetName(userService user.Service) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		if ctx.Message().Chat.ID == -571389424 {
			if ctx.Message().IsReply() {
				return handlerPaidReplyMessage(ctx)
			} else if ctx.Message().IsForwarded() {
				return handlerAddPaid(userService)(ctx)
			} else {
				return ctx.Send("Don't play", menu)
			}
		}
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
			menu.Row(btnCash),
		)
		if sendFeedback {
			sendFeedback = false
			return ctx.Send("Thank You for Feedback. Bye!!!👋")
		}
		if !sendingTask {
			if started {
				started = false
				return ctx.Reply(textHello, menu)
			} else {
				return nil
			}
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

func handlerCasher(ctx tele.Context) error {
	menu.Reply(
		menu.Row(btnCheckPaid),
	)
	return ctx.Send(textGetPaid, menu)
}

func handlerCheckPaid(ctx tele.Context) error {
	user := ctx.Get(userCtxKey).(*structs.User)
	if user.IsPaid {
		return ctx.Reply("DONE! You can learn ...", handlerInfo(ctx))
	}
	return handlerCasher(ctx)
}

func handlerLessons(bot *tele.Bot, lessonService lesson.Service, userService user.Service) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		user := ctx.Get(userCtxKey).(*structs.User)
		user, err := userService.ReadByChatID(user.ID)
		var id = user.DoingLesson

		lessons, err := lessonService.ReadByLessonID(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				menu.Reply(
					menu.Row(btnSummary),
				)
				err = ctx.Reply(congratulate, menu)
				if err != nil {
				}
				log.Println(err)
			} else {
				log.Println(err)
			}
		}

		messages, _ := parser.ParseLessonToMessage(lessons)
		for _, message := range messages {
			fmt.Printf("%+v\n%+v\n", message.Chat, *ctx.Chat())

			var recip = Recip{strconv.FormatInt(ctx.Chat().ID, 10)}

			mes, err := bot.Copy(recip, &message)
			if err != nil {
				return err
			}
			log.Printf("%+v", mes)
			if err != nil {
				log.Println("Error in Forward Lesson Messages", err)
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

		return ctx.Send(parser.ParseTimeToMessage(user.DoingLesson, user.LessonTime, lessons[0].SendDuration, false), menu)
	}
}

func handleGetLessons(lessonService lesson.Service) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		msg := ctx.Message()

		parsed, err := parser.ParseMessageToLesson(msg)
		if err != nil {
			if errors.Is(err, parser.ErrorLessonWrong) {
				return err
			} else if errors.Is(err, parser.ErrorLessonThanNeed) {
				return err
			} else {
				log.Println(err)
				return err
			}
		} else {
			err = lessonService.Store(&parsed)
			if err != nil {

				log.Println(err)
				return err
			}

			return ctx.Reply("Added new lesson")
		}
	}
}

func handleTask(userService user.Service, lessonService lesson.Service) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		user := ctx.Get(userCtxKey).(*structs.User)
		user, err := userService.ReadByChatID(user.ID)
		var id = user.DoingLesson
		lessons, err := lessonService.ReadByLessonID(id)
		if err != nil {
			log.Println(err)
		}
		startedTime := time.Now().UTC().Add(time.Hour * 5).Sub(user.LessonTime)
		duration, err := time.ParseDuration(lessons[0].SendDuration)
		if err != nil {
			log.Println(err)
		}
		if startedTime-duration < 0 {
			fmt.Printf("%.5v s\n", duration-startedTime)
			return ctx.Reply(parser.ParseTimeToMessage(user.DoingLesson, user.LessonTime, lessons[0].SendDuration, true))
		} else {
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
		var rcp tele.Recipient = Recipient{user.IsPaid}

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

func handlerPaidReplyMessage(ctx tele.Context) error {
	ctx.Message().Chat.ID = ctx.Message().ReplyTo.OriginalSender.ID
	err := ctx.Send(ctx.Message().Text)
	if err != nil {
		return err
	}

	return nil
}

func handlerAddPaid(userService user.Service) tele.HandlerFunc {
	return func(ctx tele.Context) error {
		//user := ctx.Get(userCtxKey).(*structs.User)
		err := userService.UpdatePaid(ctx.Message().OriginalSender.ID)
		if err != nil {
			log.Println(err)
			return err
		}

		return nil
	}
}

type Editable struct {
	messageID string
	chatID    int64
}

func (e Editable) MessageSig() (messageID string, chatID int64) {
	return e.messageID, e.chatID
}
