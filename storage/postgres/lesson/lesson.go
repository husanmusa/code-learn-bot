package lesson

import (
	"database/sql"
	"github.com/husanmusa/code-learn-bot/pkg/structs"
	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) StorageLesson {
	return StorageLesson{DB: db}
}

type StorageLesson struct {
	DB *sqlx.DB
}

func (s StorageLesson) Store(lesson *structs.Lesson) error {
	query := `INSERT INTO lessons (chat_id, message_id, number_of_lesson, 
                     number_of_part, type_of_part)
              VALUES ($1, $2, $3, $4, $5)
              RETURNING id`
	err := s.DB.QueryRow(
		query,
		lesson.ChatId,
		lesson.MessageId,
		lesson.NumberOfLesson,
		lesson.NumberOfPart,
		lesson.TypeOfPart,
	).Scan(&lesson.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s StorageLesson) Update(numberOfPart, typeOfPart int64, lesson *structs.Lesson) error {
	query := `UPDATE lessons
              SET
			      chat_id=$1, message_id=$2, type_of_part=$3
			  WHERE
			 	  number_of_lesson= $4 and number_of_part=$5`

	_, err := s.DB.Exec(
		query,
		lesson.ChatId,
		lesson.MessageId,
		lesson.NumberOfLesson,
		numberOfPart,
		typeOfPart,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s StorageLesson) ReadByLessonID(numberOfLesson int64) ([]structs.Lesson, error) {

	rows, err := s.DB.Query(`SELECT
		chat_id, message_id, number_of_lesson, number_of_part, type_of_part
	FROM
		lessons where number_of_lesson=$1
		`, numberOfLesson)
	if err != nil {
		return nil, err
	}

	var lessons []structs.Lesson

	for rows.Next() {
		var lesson structs.Lesson

		err = rows.Scan(
			&lesson.ChatId,
			&lesson.MessageId,
			&lesson.NumberOfLesson,
			&lesson.NumberOfPart,
			&lesson.TypeOfPart,
		)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	if len(lessons) == 0 {
		return nil, sql.ErrNoRows
	}

	return lessons, nil
}
