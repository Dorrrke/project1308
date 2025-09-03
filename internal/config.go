package internal

import "flag"

type Config struct {
	Host string
	Port int
	// TODO: DB connection string
	// TODO: Debug bool
}

func ReadConfig() Config {
	var cfg Config
	flag.StringVar(&cfg.Host, "host", "localhost", "указание адреса для запуска сервера")
	flag.IntVar(&cfg.Port, "port", 8080, "указание порта для запуска сервера")
	flag.Parse()

	return cfg
}
