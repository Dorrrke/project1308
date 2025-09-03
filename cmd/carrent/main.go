package main

import (
	"fmt"

	"github.com/Dorrrke/project1308/internal"
	"github.com/Dorrrke/project1308/internal/repository/inmemory"
)

func main() {
	fmt.Println("Server started...")
	///// TODO: конфигурация приложения
	cfg := internal.ReadConfig()
	fmt.Println(cfg)
	// TODO: конфигурация/создание хранилища\
	db := inmemory.NewInMemoryStorage()
}
