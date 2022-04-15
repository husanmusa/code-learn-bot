package user

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"strings"

	"github.com/husanmusa/code-learn-bot/pkg/structs"
)

type Storage struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) Storage {
	return Storage{DB: db}
}

func (s Storage) Store(user *structs.User) error {
	query := `INSERT INTO users (
				chat_id,
                input_name,
			   lesson_doing
              )
              VALUES ($1, $2, $3)
              RETURNING id`

	_, err := s.DB.Exec(
		query,
		user.ID,
		user.InputName,
		1,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s Storage) Update(chatID int64, user *structs.User) error {

	query := `UPDATE users
              SET
			      input_name=$1,
                  lesson_doing=$2,
                  lesson_time=$3
			  WHERE
			 	  chat_id = $4`

	_, err := s.DB.Exec(
		query,
		user.InputName,
		user.DoingLesson,
		user.LessonTime,
		chatID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s Storage) ReadByChatID(chatID int64) (*structs.User, error) {
	var (
		ld    sql.NullInt64
		lt    sql.NullTime
		user  = &structs.User{ID: chatID}
		query = `SELECT
			is_banned,
       		lesson_doing,
       		lesson_time,
			CASE WHEN input_name IS NULL THEN '' ELSE input_name END
		 FROM
			users
		 WHERE
			chat_id = $1`
	)

	err := s.DB.QueryRow(query, chatID).Scan(
		&user.IsBanned,
		&ld,
		&lt,
		&user.InputName,
	)
	if ld.Valid {
		user.DoingLesson = ld.Int64
	}
	if lt.Valid {
		user.LessonTime = lt.Time
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s Storage) ReadMany() ([]*structs.User, error) {
	var (
		query strings.Builder
	)

	query.WriteString(`SELECT
		users.chat_id,
		users.is_banned,
		users.lesson_doing,
		CASE WHEN users.input_name IS NULL THEN '' ELSE users.input_name END
	FROM
		users
		`)

	rows, err := s.DB.Query(query.String())
	if err != nil {
		return nil, err
	}

	var users []*structs.User

	for rows.Next() {
		var user structs.User

		err = rows.Scan(
			&user.ID,
			&user.IsBanned,
			&user.DoingLesson,
			&user.InputName,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return users, nil
}

func (Storage) IsErrNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
