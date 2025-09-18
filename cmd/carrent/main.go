package main

import (
	"github.com/Dorrrke/project1308/internal"
	"github.com/Dorrrke/project1308/internal/repository/db"
	"github.com/Dorrrke/project1308/internal/server"
	"github.com/Dorrrke/project1308/pkg/logger"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := internal.ReadConfig()
	log := logger.Init(cfg.Debug)
	log.Debug().Any("config", cfg).Send()

	log.Info().Msg("Server starting...")
	// TODO: конфигурация/создание хранилища\
	// database := inmemory.NewInMemoryStorage()
	database, err := db.NewStorage(cfg.DSN)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create db")
	}

	if err := db.Migrations(cfg.DSN, cfg.MigratePath); err != nil {
		log.Fatal().Err(err).Msg("failed to migrate db")
	}

	srv := server.NewServer(cfg, database, &log)

	if err := srv.Run(); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
