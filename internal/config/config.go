package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	defaultValueServerAddress   = ":8080"
	defaultValueURLAddress      = "http://localhost:8080"
	defaultValueLogLevel        = "debug"
	defaultValueFileStoragePath = ""
	defaultValueDatabaseDSN     = ""
	defaultValueAuditFile       = ""
	defaultValueAuditURL        = ""
	defaultValueEnableHTTPS     = false
	defaultValueAutoCert        = false
	defaultValueTLSCerfFile     = "./resources/ssl/cert.pem"
	defaultValueTLSKeyFile      = "./resources/ssl/key.pem"
	defaultValueJSONConfigFile  = ""
)

type Config struct {
	ServerAddress   string `json:"server_address,omitempty"`
	URLAddress      string `json:"base_url,omitempty"`
	LogLevel        string `json:"log_level,omitempty"`
	FileStoragePath string `json:"file_storage_path,omitempty"`
	DatabaseDSN     string `json:"database_dsn,omitempty"`
	AuditFile       string `json:"audit_file,omitempty"`
	AuditURL        string `json:"audit_url,omitempty"`
	EnableHTTPS     bool   `json:"enable_https,omitempty"`
	AutoCert        bool   `json:"auto_cert,omitempty"`
	TLSCerfFile     string `json:"tls_cerf_file,omitempty"`
	TLSKeyFile      string `json:"tls_key_file,omitempty"`
	jsonConfigFile  string
}

var Values Config

func Init() error {
	cfg := Config{}
	flag.StringVar(&cfg.ServerAddress, "a", defaultValueServerAddress, "address of HTTP server to start")
	flag.StringVar(&cfg.URLAddress, "b", defaultValueURLAddress, "server address in shortened URLs")
	flag.StringVar(&cfg.LogLevel, "l", defaultValueLogLevel, "log level")
	flag.StringVar(&cfg.FileStoragePath, "f", defaultValueFileStoragePath, "path to file with short urls")
	flag.StringVar(&cfg.DatabaseDSN, "d", defaultValueDatabaseDSN, "connection settings for pgsql") // postgres://user:pass@localhost:5432/test-db?sslmode=disable
	flag.StringVar(&cfg.AuditFile, "audit-file", defaultValueAuditFile, "path to audit file")
	flag.StringVar(&cfg.AuditURL, "audit-url", defaultValueAuditURL, "url to submit audit data")
	flag.BoolVar(&cfg.EnableHTTPS, "s", defaultValueEnableHTTPS, "enable https")
	flag.BoolVar(&cfg.AutoCert, "auto-cert", defaultValueAutoCert, "generate TLS certificate automatically (this option overrides -cert-file and -key-file options)")
	flag.StringVar(&cfg.TLSCerfFile, "cert-file", defaultValueTLSCerfFile, "path to TLS certificate")
	flag.StringVar(&cfg.TLSKeyFile, "key-file", defaultValueTLSKeyFile, "path to TLS key")

	flag.StringVar(&cfg.jsonConfigFile, "config", defaultValueJSONConfigFile, "path to JSON config file")
	flag.StringVar(&cfg.jsonConfigFile, "c", defaultValueJSONConfigFile, "path to JSON config file (shorthand)")

	flag.Parse()

	jsonConfigFile, ok := os.LookupEnv("CONFIG")
	if ok {
		cfg.jsonConfigFile = jsonConfigFile
	}

	if cfg.jsonConfigFile != defaultValueJSONConfigFile {
		jsonConfigData, err := os.ReadFile(cfg.jsonConfigFile)
		if err != nil {
			return fmt.Errorf("could not scan JSON config file %s: %w", cfg.jsonConfigFile, err)
		}
		var jsonConfigValues Config
		err = json.Unmarshal(jsonConfigData, &jsonConfigValues)
		if err != nil {
			return fmt.Errorf("could not unmarshall JSON config file %s: %w", cfg.jsonConfigFile, err)
		}

		if cfg.ServerAddress == defaultValueServerAddress && jsonConfigValues.ServerAddress != "" {
			cfg.ServerAddress = jsonConfigValues.ServerAddress
		}
		if cfg.URLAddress == defaultValueURLAddress && jsonConfigValues.URLAddress != "" {
			cfg.URLAddress = jsonConfigValues.URLAddress
		}
		if cfg.LogLevel == defaultValueLogLevel && jsonConfigValues.LogLevel != "" {
			cfg.LogLevel = jsonConfigValues.LogLevel
		}
		if cfg.FileStoragePath == defaultValueFileStoragePath && jsonConfigValues.FileStoragePath != "" {
			cfg.FileStoragePath = jsonConfigValues.FileStoragePath
		}
		if cfg.DatabaseDSN == defaultValueDatabaseDSN && jsonConfigValues.DatabaseDSN != "" {
			cfg.DatabaseDSN = jsonConfigValues.DatabaseDSN
		}
		if cfg.AuditFile == defaultValueAuditFile && jsonConfigValues.AuditFile != "" {
			cfg.AuditFile = jsonConfigValues.AuditFile
		}
		if cfg.AuditURL == defaultValueAuditURL && jsonConfigValues.AuditURL != "" {
			cfg.AuditURL = jsonConfigValues.AuditURL
		}
		if cfg.EnableHTTPS == defaultValueEnableHTTPS && jsonConfigValues.EnableHTTPS != false {
			cfg.EnableHTTPS = jsonConfigValues.EnableHTTPS
		}
		if cfg.AutoCert == defaultValueAutoCert && jsonConfigValues.AutoCert != false {
			cfg.AutoCert = jsonConfigValues.AutoCert
		}
		if cfg.TLSCerfFile == defaultValueTLSCerfFile && jsonConfigValues.TLSCerfFile != "" {
			cfg.TLSCerfFile = jsonConfigValues.TLSCerfFile
		}
		if cfg.TLSKeyFile == defaultValueTLSKeyFile && jsonConfigValues.TLSKeyFile != "" {
			cfg.TLSKeyFile = jsonConfigValues.TLSKeyFile
		}
	}

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
