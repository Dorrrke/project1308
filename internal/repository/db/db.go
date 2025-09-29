package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
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

func (s *Storage) Close(ctx context.Context) error {
	return s.userStorage.db.Close(ctx)
}

func Migrations(dsn string, migratePath string, log *zerolog.Logger) error { // /migrations
	mPath := fmt.Sprintf("file://%s", migratePath) // file:///migrations

	m, err := migrate.New(mPath, dsn)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
		log.Debug().Msg("No changes in migrations")
	}

	log.Debug().Msg("Migrations completed")
	return nil
}
