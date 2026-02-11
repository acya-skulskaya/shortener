package config

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"dario.cat/mergo"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServerAddress   string `json:"server_address,omitempty" env:"SERVER_ADDRESS" env-default:":8080"`
	URLAddress      string `json:"base_url,omitempty" env:"BASE_URL" env-default:"http://localhost:8080"`
	LogLevel        string `json:"log_level,omitempty" env:"LOG_LEVEL" env-default:"debug"`
	FileStoragePath string `json:"file_storage_path,omitempty" env:"FILE_STORAGE_PATH" env-default:""`
	DatabaseDSN     string `json:"database_dsn,omitempty" env:"DATABASE_DSN" env-default:""`
	AuditFile       string `json:"audit_file,omitempty" env:"AUDIT_FILE" env-default:""`
	AuditURL        string `json:"audit_url,omitempty" env:"AUDIT_URL" env-default:"false"`
	EnableHTTPS     bool   `json:"enable_https,omitempty" env:"ENABLE_HTTPS" env-default:"false"`
	AutoCert        bool   `json:"auto_cert,omitempty" env:"AUTO_CERT" env-default:"false"`
	TLSCerfFile     string `json:"tls_cerf_file,omitempty" env:"TLS_CERT_FILE" env-default:"./resources/ssl/cert.pem"`
	TLSKeyFile      string `json:"tls_key_file,omitempty" env:"TLS_KEY_FILE" env-default:"./resources/ssl/key.pem"`
	jsonConfigFile  string
}

var Values Config

func Init() error {
	cfgFlag := Config{}
	flag.StringVar(&cfgFlag.ServerAddress, "a", "", "address of HTTP server to start")
	flag.StringVar(&cfgFlag.URLAddress, "b", "", "server address in shortened URLs")
	flag.StringVar(&cfgFlag.LogLevel, "l", "", "log level")
	flag.StringVar(&cfgFlag.FileStoragePath, "f", "", "path to file with short urls")
	flag.StringVar(&cfgFlag.DatabaseDSN, "d", "", "connection settings for pgsql") // postgres://user:pass@localhost:5432/test-db?sslmode=disable
	flag.StringVar(&cfgFlag.AuditFile, "audit-file", "", "path to audit file")
	flag.StringVar(&cfgFlag.AuditURL, "audit-url", "", "url to submit audit data")
	flag.BoolVar(&cfgFlag.EnableHTTPS, "s", false, "enable https")
	flag.BoolVar(&cfgFlag.AutoCert, "auto-cert", false, "generate TLS certificate automatically (this option overrides -cert-file and -key-file options)")
	flag.StringVar(&cfgFlag.TLSCerfFile, "cert-file", "", "path to TLS certificate")
	flag.StringVar(&cfgFlag.TLSKeyFile, "key-file", "", "path to TLS key")

	flag.StringVar(&cfgFlag.jsonConfigFile, "config", "", "path to JSON config file")
	flag.StringVar(&cfgFlag.jsonConfigFile, "c", "", "path to JSON config file (shorthand)")

	flag.Parse()

	jsonConfigFileEnv, ok := os.LookupEnv("CONFIG")
	if ok {
		cfgFlag.jsonConfigFile = jsonConfigFileEnv
	}

	var cfg Config
	if cfgFlag.jsonConfigFile != "" {
		err := cleanenv.ReadConfig(cfgFlag.jsonConfigFile, &cfg)
		if err != nil {
			return fmt.Errorf("could not scan JSON config file or env variables %s: %w", cfgFlag.jsonConfigFile, err)
		}
	} else {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return fmt.Errorf("could not scan env variables: %w", err)
		}
	}

	err := mergo.Merge(&cfg, cfgFlag, mergo.WithOverride)
	if err != nil {
		return fmt.Errorf("could not merge config sources: %w", err)
	}

	if cfg.AutoCert {
		cfg.TLSCerfFile = ""
		cfg.TLSKeyFile = ""
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
