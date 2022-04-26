package structs

import "time"

type User struct {
	ID          int64
	InputName   string
	Firstname   string
	LastName    string
	DoingLesson int64
	LessonTime  time.Time
	IsBanned    bool
	IsPaid      bool
}

type Lesson struct {
	ID             int64
	ChatId         int64
	MessageId      int
	NumberOfLesson int64
	NumberOfPart   int64
	SendDuration   string
}
