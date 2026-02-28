package grpc

import (
	"crypto/tls"
	"fmt"

	"github.com/acya-skulskaya/shortener/internal/config"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc/credentials"
)

func LoadTLSCredentials() (credentials.TransportCredentials, error) {
	if !config.Values.EnableHTTPS {
		return nil, fmt.Errorf("TLS was not enabled in configuration")
	}

	if config.Values.AutoCert {
		manager := autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache("golang-autocert"),
		}

		return credentials.NewTLS(manager.TLSConfig()), nil
	}

	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(config.Values.TLSCerfFile, config.Values.TLSKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not load server certificate: %w", err)
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}
