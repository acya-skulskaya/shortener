package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress   string
	URLAddress      string
	LogLevel        string
	FileStoragePath string
	DatabaseDSN     string
	AuditFile       string
	AuditURL        string
}

var Values Config

func Init() {
	cfg := Config{}
	flag.StringVar(&cfg.ServerAddress, "a", ":8080", "address of HTTP server to start")
	flag.StringVar(&cfg.URLAddress, "b", "http://localhost:8080", "server address in shortened URLs")
	flag.StringVar(&cfg.LogLevel, "l", "debug", "log level")
	flag.StringVar(&cfg.FileStoragePath, "f", "", "path to file with short urls")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "connection settings for pgsql") // postgres://user:pass@localhost:5432/test-db?sslmode=disable
	flag.StringVar(&cfg.AuditFile, "audit-file", "", "path to audit file")
	flag.StringVar(&cfg.AuditURL, "audit-url", "", "url to submit audit data")

	flag.Parse()

	serverAddress, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok {
		cfg.ServerAddress = serverAddress
	}

	baseURL, ok := os.LookupEnv("BASE_URL")
	if ok {
		cfg.URLAddress = baseURL
	}

	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if ok {
		cfg.LogLevel = logLevel
	}

	fileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		cfg.FileStoragePath = fileStoragePath
	}

	databaseDSN, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		cfg.DatabaseDSN = databaseDSN
	}

	auditFile, ok := os.LookupEnv("AUDIT_FILE")
	if ok {
		cfg.AuditFile = auditFile
	}

	auditURL, ok := os.LookupEnv("AUDIT_URL")
	if ok {
		cfg.AuditURL = auditURL
	}

	Values = cfg
}
