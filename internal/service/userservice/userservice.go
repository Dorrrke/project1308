package userservice

import (
	erros "github.com/Dorrrke/project1308/internal/domain/user/errors"
	"github.com/Dorrrke/project1308/internal/domain/user/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
)

type UserStorage interface {
	SaveUser(user models.User) error
	GetUser(userReq models.UserRequest) (models.User, error)
	GetUserByID(uid string) (models.User, error)
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
	if err := us.valid.Struct(user); err != nil {
		return err
	}

	uid := uuid.New().String()
	user.UID = uid

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	return us.db.SaveUser(user)
}

func (us *UserService) LoginUser(userReq models.UserRequest) (models.User, error) {
	dbUser, err := us.db.GetUser(userReq)
	if err != nil {
		return models.User{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(userReq.Password)); err != nil {
		return models.User{}, erros.ErrInvalidPassword
	}

	return dbUser, nil
}

func (us *UserService) GetUserByID(uid string) (models.User, error) {
	return us.db.GetUserByID(uid)
}
