package main

import (
	"fmt"

	"github.com/Dorrrke/project1308/internal"
	"github.com/Dorrrke/project1308/internal/repository/inmemory"
	"github.com/Dorrrke/project1308/internal/server"
)

func main() {
	fmt.Println("Server started...")
	///// TODO: конфигурация приложения
	cfg := internal.ReadConfig()
	fmt.Println(cfg)
	// TODO: конфигурация/создание хранилища\
	db := inmemory.NewInMemoryStorage()

	// TODO: конфигурация и запуск веб-сервера
	srv := server.NewServer(cfg, db)

	if err := srv.Run(); err != nil {
		panic(err)
	}
}
