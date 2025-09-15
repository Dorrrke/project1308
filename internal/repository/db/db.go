package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
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

func Migrations(dsn string, migratePath string) error { // /migrations
	mPath := fmt.Sprintf("file://%s", migratePath) // file:///migrations

	m, err := migrate.New(mPath, dsn)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
		log.Println("Database already up to date")
	}

	log.Println("Migration completed")
	return nil
}
