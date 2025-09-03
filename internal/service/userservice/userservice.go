package userservice

import (
	"fmt"

	erros "github.com/Dorrrke/project1308/internal/domain/user/errors"
	"github.com/Dorrrke/project1308/internal/domain/user/models"

	"github.com/go-playground/validator/v10"
)

type UserStorage interface {
	SaveUser(user models.User) error
	GetUser(userReq models.UserRequest) (models.User, error)
}

type UserService struct {
	db    UserStorage
	valid *validator.Validate
}

func NewUserService(db UserStorage) *UserService {
	return &UserService{
		db:    db,
		valid: validator.New(),
	}
}

func (us *UserService) SaveUser(user models.User) error {
	// TODO: валидация
	if err := us.valid.Struct(user); err != nil {
		return err
	}

	fmt.Println()
	// TODO: хеширование пароля

	return us.db.SaveUser(user)
}

func (us *UserService) LoginUser(userReq models.UserRequest) (models.User, error) {
	dbUser, err := us.db.GetUser(userReq)
	if err != nil {
		return models.User{}, err
	}

	// TODO: проверка хеша и пароля
	if dbUser.Password != userReq.Password {
		return models.User{}, erros.ErrInvalidPassword
	}

	return dbUser, nil
}
