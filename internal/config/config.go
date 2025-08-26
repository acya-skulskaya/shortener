package config

import "flag"

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

	Values = cfg
}
