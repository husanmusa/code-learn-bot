package parser

import (
	"errors"
	"fmt"
	"github.com/husanmusa/code-learn-bot/pkg/structs"
	tele "gopkg.in/telebot.v3"
	"strconv"
	"strings"
	"time"
)

func ParseMessageToLesson(msg *tele.Message) (structs.Lesson, error) {
	var (
		lesson       structs.Lesson
		text         []string
		sendDuration string
	)
	if msg.Text != "" {
		text = strings.Split(msg.Text, "\n")
	} else {
		text = strings.Split(msg.Caption, "\n")
	}
	if len(text[0]) > 7 && len(text[1]) > 5 {
		if text[0][:7] == "Lesson-" && text[1][:5] == "Part-" {
			numberOfLesson, err := strconv.Atoi(text[0][7:len(text[0])])
			if err != nil {
				return structs.Lesson{}, err
			}
			numberOfPart, err := strconv.Atoi(text[1][5:len(text[1])])
			if err != nil {
				return structs.Lesson{}, err
			}
			if numberOfPart == 1 {
				sendDuration = text[2][5:len(text[2])]
				if err != nil {
					return structs.Lesson{}, err
				}
			}
			lesson.ChatId = msg.Chat.ID
			lesson.MessageId = msg.ID
			lesson.NumberOfLesson = int64(numberOfLesson)
			lesson.NumberOfPart = int64(numberOfPart)
			if numberOfPart == 1 {
				lesson.SendDuration = sendDuration
			}

			return lesson, nil
		} else {
			return structs.Lesson{}, ErrorLessonThanNeed
		}
	} else {
		return structs.Lesson{}, ErrorLessonWrong
	}
}

func ParseLessonToMessage(lessons []structs.Lesson) ([]tele.Message, error) {
	var messages []tele.Message

	for _, lesson := range lessons {
		var chat tele.Chat
		chat.ID = lesson.ChatId
		var message tele.Message
		message.ID = lesson.MessageId
		message.Chat = &chat
		messages = append(messages, message)
	}
	return messages, nil
}

func ParseTimeToMessage(doing int64, timer time.Time, timeDuration string, check bool) string {
	timer.Add(time.Hour * 5)
	hour, minute := ParseStringToInt(timeDuration)
	if check {
		leftHour := time.Now().Hour() - timer.Hour()
		leftMinute := time.Now().Hour() - timer.Hour()
		//leftSecond := time.Now().Second() - timer.Second()
		resp := fmt.Sprintf("❌ Таймер еще не вышел, вы сможете загрузить результат через(h:m) %02d:%02d",
			leftHour, leftMinute)
		if leftMinute == 0 {
			resp = fmt.Sprintf("❌ Таймер еще не вышел, вы сможете загрузить результат через менее чем через 1 минуту")
		}
		return resp
	}
	return fmt.Sprintf("⏱ Количество minut на выполнение задачи: %d\n\n📆 Добавить результат вы сможете после: %02d-%02d %02d:%02d",
		doing, timer.Month(), timer.Day(), timer.Hour()+hour, timer.Minute()+minute)
}

func ParseStringToInt(t string) (int, int) {
	if len(t) == 4 {
		return int(t[0]) - 48, int(t[2]) - 48
	} else if t[1] == 'm' {
		return 0, int(t[0]) - 48
	} else if t[1] == 'h' {
		return int(t[0]) - 48, 0
	}
	return -1, -1
}

/*
Lesson-1
Part-1
Type-text
*/

var ErrorLessonThanNeed = errors.New("wrong lesson creating")
var ErrorLessonWrong = errors.New("it is not lesson create post")
