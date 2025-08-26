package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress string
	URLAddress    string
}

var Values Config

func Init() {
	cfg := Config{}
	flag.StringVar(&cfg.ServerAddress, "a", ":8080", "address of HTTP server to start")
	flag.StringVar(&cfg.URLAddress, "b", "http://localhost:8080", "server address in shortened URLs")

	flag.Parse()

	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		cfg.ServerAddress = serverAddress
	}

	if baseUrl := os.Getenv("BASE_URL"); baseUrl != "" {
		cfg.URLAddress = baseUrl
	}

	Values = cfg
}
