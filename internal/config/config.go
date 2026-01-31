package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

// generate:reset
type Config struct {
	ServerAddress   string
	URLAddress      string
	LogLevel        string
	FileStoragePath string
	DatabaseDSN     string
	AuditFile       string
	AuditURL        string
	EnableHTTPS     bool
	AutoCert        bool
	TLSCerfFile     string
	TLSKeyFile      string
}

var Values Config

func Init() error {
	cfg := Config{}
	flag.StringVar(&cfg.ServerAddress, "a", ":8080", "address of HTTP server to start")
	flag.StringVar(&cfg.URLAddress, "b", "http://localhost:8080", "server address in shortened URLs")
	flag.StringVar(&cfg.LogLevel, "l", "debug", "log level")
	flag.StringVar(&cfg.FileStoragePath, "f", "", "path to file with short urls")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "connection settings for pgsql") // postgres://user:pass@localhost:5432/test-db?sslmode=disable
	flag.StringVar(&cfg.AuditFile, "audit-file", "", "path to audit file")
	flag.StringVar(&cfg.AuditURL, "audit-url", "", "url to submit audit data")
	flag.BoolVar(&cfg.EnableHTTPS, "s", false, "enable https")
	flag.BoolVar(&cfg.AutoCert, "auto-cert", false, "generate TLS certificate automatically (this option overrides -cert-file and -key-file options)")
	flag.StringVar(&cfg.TLSCerfFile, "cert-file", "./resources/ssl/cert.pem", "path to TLS certificate")
	flag.StringVar(&cfg.TLSKeyFile, "key-file", "./resources/ssl/key.pem", "path to TLS key")

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

	enableHTTPS, ok := os.LookupEnv("ENABLE_HTTPS")
	if ok {
		enableHTTPS = strings.ToLower(enableHTTPS)
		if enableHTTPS == "true" || enableHTTPS == "1" {
			cfg.EnableHTTPS = true
		}
	}

	certFile, ok := os.LookupEnv("TLS_CERT_FILE")
	if ok {
		cfg.TLSCerfFile = certFile
	}

	keyFile, ok := os.LookupEnv("TLS_KEY_FILE")
	if ok {
		cfg.TLSKeyFile = keyFile
	}

	autoCert, ok := os.LookupEnv("ENABLE_HTTPS")
	if ok {
		autoCert = strings.ToLower(autoCert)
		if autoCert == "true" || autoCert == "1" {
			cfg.AutoCert = true
			cfg.TLSCerfFile = ""
			cfg.TLSKeyFile = ""
		}
	}

	if cfg.EnableHTTPS && !cfg.AutoCert {
		if _, err := os.Stat(cfg.TLSCerfFile); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("could not find TLS certificate file at %s: %w", cfg.TLSCerfFile, err)
		}
		if _, err := os.Stat(cfg.TLSKeyFile); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("could not find TLS key file at %s: %w", cfg.TLSKeyFile, err)
		}
	}

	Values = cfg

	return nil
}
