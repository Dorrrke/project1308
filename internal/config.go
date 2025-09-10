package internal

import "flag"

type Config struct {
	Host string
	Port int
	DSN  string
	// TODO: Debug bool
}

func ReadConfig() Config {
	var cfg Config
	flag.StringVar(&cfg.Host, "host", "localhost", "указание адреса для запуска сервера")
	flag.IntVar(&cfg.Port, "port", 8080, "указание порта для запуска сервера")
	flag.StringVar(&cfg.DSN, "dsn", "postgres://postgres:postgres@0.0.0.0:5432/postgres", "указание строки подключения к БД")
	flag.Parse()

	return cfg
}
