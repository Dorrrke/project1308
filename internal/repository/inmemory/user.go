package inmemory

import (
	userErrors "github.com/Dorrrke/project1308/internal/domain/user/errors"
	userDomain "github.com/Dorrrke/project1308/internal/domain/user/models"

	"github.com/google/uuid"
)

func (storage *Storage) SaveUser(user userDomain.User) error {
	for _, us := range storage.users {
		if user.Email == us.Email || user.Phone == us.Phone {
			return userErrors.ErrUserAlreadyExists
		}
	}
	uid := uuid.New().String()
	_, ok := storage.users[uid]

	if ok {
		uid = uuid.New().String()
	}

	user.UID = uid
	storage.users[user.UID] = user
	return nil
}

func (storage *Storage) GetUser(userReq userDomain.UserRequest) (userDomain.User, error) {
	for _, us := range storage.users {
		if userReq.Email == us.Email {
			return us, nil
		}
	}
	return userDomain.User{}, userErrors.ErrUserNoExists
}
