package storage

import "github.com/husanmusa/code-learn-bot/pkg/structs"

type LessonStorage interface {
	Store(*structs.Lesson) error
	Update(numberOfPart int64, lesson *structs.Lesson) error
	ReadByLessonID(int64) ([]structs.Lesson, error)
}
