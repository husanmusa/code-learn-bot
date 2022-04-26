package storage

import (
	"github.com/husanmusa/code-learn-bot/pkg/structs"
)

type UserStorage interface {
	Store(*structs.User) error
	Update(chatID int64, user *structs.User) error
	UpdatePaid(chatID int64) error
	ReadByChatID(int64) (*structs.User, error)
	ReadMany() ([]*structs.User, error)
	IsErrNotFound(error) bool
}
