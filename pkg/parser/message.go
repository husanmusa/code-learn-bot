package parser

import (
	"fmt"
	"github.com/husanmusa/code-learn-bot/pkg/structs"
	tele "gopkg.in/telebot.v3"
	"strconv"
	"strings"
	"time"
)

func ParseMessageToLesson(msg *tele.Message) (structs.Lesson, error) {
	var lesson structs.Lesson
	text := strings.Split(msg.Text, "\n")
	numberOfLesson, err := strconv.Atoi(text[0][7:len(text[0])])
	if err != nil {
		return structs.Lesson{}, nil
	}
	numberOfPart, err := strconv.Atoi(text[1][5:len(text[1])])
	if err != nil {
		return structs.Lesson{}, nil
	}
	typeOfPart := text[2][5:len(text[2])]
	if err != nil {
		return structs.Lesson{}, nil
	}
	lesson.ChatId = msg.Chat.ID
	lesson.MessageId = msg.ID
	lesson.NumberOfLesson = int64(numberOfLesson)
	lesson.NumberOfPart = int64(numberOfPart)
	lesson.TypeOfPart = typeOfPart

	return lesson, nil
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

func ParseTimeToMessage(timer time.Time) string {
	timer.Add(time.Hour * 5)
	return fmt.Sprintf("⏱ Количество minut на выполнение задачи: 3\n\n📆 Добавить результат вы сможете после: %02d-%02d %02d:%02d",
		timer.Month(), timer.Day(), timer.Hour(), timer.Minute()+1)
}

/*
Lesson-1
Part-1
Type-text
*/
