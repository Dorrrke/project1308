package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Dorrrke/project1308/internal"
	"github.com/Dorrrke/project1308/internal/repository/db"
	"github.com/Dorrrke/project1308/internal/server"
	"github.com/Dorrrke/project1308/pkg/logger"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

		<-c
		cancel()
	}()

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

	if err = db.Migrations(cfg.DSN, cfg.MigratePath, &log); err != nil {
		log.Fatal().Err(err).Msg("failed to migrate db")
	}

	srv := server.NewServer(cfg, database, &log)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = srv.Run(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	wg.Add(1)
	go func() {
		<-ctx.Done()
		log.Debug().Msg("Shutdown signal received")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("failed to shutdown server")
		}
		if err := database.Close(ctx); err != nil {
			log.Error().Err(err).Msg("failed to shutdown database")
		}

		log.Info().Msg("Server stopped")
	}()

	wg.Wait()
}
