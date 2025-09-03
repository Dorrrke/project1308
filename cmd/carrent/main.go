package main

import (
	"fmt"

	"github.com/Dorrrke/project1308/internal"
)

func main() {
	fmt.Println("Server started...")
	///// TODO: конфигурация приложения
	cfg := internal.ReadConfig()
	fmt.Println(cfg)
	// TODO: конфигурация/создание хранилища
}
