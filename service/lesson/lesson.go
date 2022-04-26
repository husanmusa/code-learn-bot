package lesson

import (
	"github.com/husanmusa/code-learn-bot/storage"
	"log"

	"github.com/husanmusa/code-learn-bot/pkg/structs"
)

type Service struct {
	LessonStorage storage.LessonStorage
}

func NewService(s storage.LessonStorage) Service {
	if s == nil {
		log.Fatal("storage is nil")
	}
	return Service{LessonStorage: s}
}

func (s Service) Store(lesson *structs.Lesson) error {

	err := s.LessonStorage.Store(lesson)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) ReadByLessonID(numberOfLesson int64) ([]structs.Lesson, error) {
	return s.LessonStorage.ReadByLessonID(numberOfLesson)
}

func (s Service) Update(numberOfPart int64, lesson *structs.Lesson) error {
	return s.LessonStorage.Update(numberOfPart, lesson)
}
