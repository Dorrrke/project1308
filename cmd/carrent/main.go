package main

import (
	"fmt"
	"log"

	"github.com/Dorrrke/project1308/internal"
	"github.com/Dorrrke/project1308/internal/repository/db"
	"github.com/Dorrrke/project1308/internal/server"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	fmt.Println("Server started...")
	///// TODO: конфигурация приложения
	cfg := internal.ReadConfig()
	fmt.Println(cfg)
	// TODO: конфигурация/создание хранилища\
	// database := inmemory.NewInMemoryStorage()
	database, err := db.NewStorage(cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Migrations(cfg.DSN, cfg.MigratePath); err != nil {
		log.Fatal(err)
	}
	// TODO: конфигурация и запуск веб-сервера
	srv := server.NewServer(cfg, database)

	if err := srv.Run(); err != nil {
		panic(err)
	}
}
