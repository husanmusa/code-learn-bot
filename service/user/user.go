package user

import (
	"errors"
	"fmt"
	"github.com/husanmusa/code-learn-bot/storage"
	"log"

	"github.com/husanmusa/code-learn-bot/pkg/structs"
)

var _errUserBanned = errors.New("chat is banned")

type Service struct {
	UserStorage storage.UserStorage
}

func NewService(s storage.UserStorage) Service {
	if s == nil {
		log.Fatal("storage is nil")
	}
	return Service{UserStorage: s}
}

func (s Service) Auth(chatID int64) (*structs.User, error) {
	// check. does have user in db?
	existingUser, err := s.UserStorage.ReadByChatID(chatID)
	if err != nil {
		if s.UserStorage.IsErrNotFound(err) { //if user not found
			user := structs.User{
				ID: chatID,
			} //user is not found, then store him.
			if err = s.UserStorage.Store(&user); err != nil {
				return nil, err
			}
			return &user, nil
		} else {
			return nil, err
		}
	}

	if existingUser.IsBanned {
		return nil, fmt.Errorf("%d %w", existingUser.ID, _errUserBanned)
	}

	return existingUser, nil
}

func (s Service) ReadMany() ([]*structs.User, error) {
	return s.UserStorage.ReadMany()
}

func (s Service) ReadByChatID(chatID int64) (*structs.User, error) {
	return s.UserStorage.ReadByChatID(chatID)
}

func (s Service) Update(chatID int64, user *structs.User) error {
	return s.UserStorage.Update(chatID, user)
}

func (s Service) IsErrUserBanned(err error) bool {
	return errors.Is(errors.Unwrap(err), _errUserBanned)
}
