package internal

import (
	"cmp"
	"flag"
	"os"
)

const (
	defDNS         = "postgres://postgres:postgres@0.0.0.0:5432/postgres?sslmode=disable"
	defHost        = "0.0.0.0"
	defPort        = 8080
	defMigratePath = "migrations"
)

type Config struct {
	Host        string
	Port        int
	DSN         string
	MigratePath string
}

func ReadConfig() Config {
	var cfg Config
	flag.StringVar(&cfg.Host, "host", defHost, "указание адреса для запуска сервера")
	flag.IntVar(&cfg.Port, "port", defPort, "указание порта для запуска сервера")
	flag.StringVar(&cfg.DSN, "dsn", defDNS, "указание строки подключения к БД")
	flag.StringVar(&cfg.MigratePath, "migrate-path", defMigratePath, "указание пути к миграциям")
	flag.Parse()

	cfg.DSN = cmp.Or(os.Getenv("DB_DNS"), defDNS)
	cfg.MigratePath = cmp.Or(os.Getenv("MIGRATE_PATH"), defMigratePath)

	return cfg
}
