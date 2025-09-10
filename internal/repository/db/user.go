package db

import (
	"context"
	"time"

	userDomain "github.com/Dorrrke/project1308/internal/domain/user/models"

	"github.com/jackc/pgx/v5"
)

type userStorage struct {
	db *pgx.Conn
}

func (us *userStorage) SaveUser(user userDomain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := us.db.Exec(ctx, "INSERT INTO users (uid, name, email, password, age, phone) VALUES ($1, $2, $3, $4, $5, $6)",
		user.UID, user.Name, user.Email, user.Password, user.Age, user.Phone)
	if err != nil {
		// TODO: обработка ошибки при существующем uid
		// TODO: обработка ошибки - пользователь уже есть
		return err
	}

	return nil
}

func (us *userStorage) GetUser(userReq userDomain.UserRequest) (userDomain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user userDomain.User
	err := us.db.QueryRow(ctx, "SELECT * FROM users WHERE email = $1", userReq.Email).
		Scan(&user.UID, &user.Name, &user.Email, &user.Password, &user.Age, &user.Phone)
	if err != nil {
		return userDomain.User{}, err
	}

	return user, nil
}

func (us *userStorage) GetUserByID(uid string) (userDomain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user userDomain.User
	err := us.db.QueryRow(ctx, "SELECT * FROM users WHERE uid = $1", uid).
		Scan(&user.UID, &user.Name, &user.Email, &user.Password, &user.Age, &user.Phone)
	if err != nil {
		return userDomain.User{}, err
	}

	return user, nil
}
