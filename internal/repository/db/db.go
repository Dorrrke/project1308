package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	userStorage
	carsStorage
}

func NewStorage(connStr string) (*Storage, error) {
	db, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	return &Storage{
		carsStorage: carsStorage{db: db},
		userStorage: userStorage{db: db},
	}, nil
}
