package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/acya-skulskaya/shortener/api/shortener"
	"github.com/acya-skulskaya/shortener/internal/config"
	grpcHandler "github.com/acya-skulskaya/shortener/internal/handler/grpc"
	httpHandler "github.com/acya-skulskaya/shortener/internal/handler/http"
	"github.com/acya-skulskaya/shortener/internal/interceptors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/acya-skulskaya/shortener/internal/observer/audit/publisher"
	"github.com/acya-skulskaya/shortener/internal/observer/audit/subscribers"
	interfaces "github.com/acya-skulskaya/shortener/internal/repository/interface"
	shorturlindb "github.com/acya-skulskaya/shortener/internal/repository/short_url_in_db"
	shorturlinmemory "github.com/acya-skulskaya/shortener/internal/repository/short_url_in_memory"
	shorturljsonfile "github.com/acya-skulskaya/shortener/internal/repository/short_url_json_file"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
)

const (
	shutdownPeriod = 20 * time.Second
)

func main() {
	printBuildInfo()

	if err := config.Init(); err != nil {
		log.Fatalf("could not init configuration: %v", err)
	}

	if err := logger.Init(config.Values.LogLevel); err != nil {
		log.Fatalf("failed init logger: %v", err)
	}

	// Setup signal context
	rootCtx, stopRootCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stopRootCtx()

	auditPublisher := publisher.NewAuditPublisher()
	if config.Values.AuditFile != "" || config.Values.AuditURL != "" {
		if config.Values.AuditFile != "" {
			fileAuditSubscriber, err := subscribers.NewFileAuditSubscriber(rootCtx, config.Values.AuditFile)
			if err != nil {
				logger.Log.Fatal("could not create file audit subscriber", zap.Error(err))
			} else {
				auditPublisher.Subscribe(fileAuditSubscriber)
			}
		}
		if config.Values.AuditURL != "" {
			httpAuditSubscriber := subscribers.NewHTTPAuditSubscriber(rootCtx, config.Values.AuditURL)
			auditPublisher.Subscribe(httpAuditSubscriber)
		}
	}

	var repo interfaces.ShortURLRepository
	if len(config.Values.DatabaseDSN) != 0 {
		repo, err := shorturlindb.NewInDBShortURLRepository(config.Values.DatabaseDSN)
		if err != nil {
			logger.Log.Fatal("failed to init db storage", zap.Error(err))
		}
		defer repo.Close()
		logger.Log.Info("using db storage", zap.String("DatabaseDSN", config.Values.DatabaseDSN))
	} else if len(config.Values.FileStoragePath) != 0 {
		repo = &shorturljsonfile.JSONFileShortURLRepository{FileStoragePath: config.Values.FileStoragePath}
		logger.Log.Info("using file storage", zap.String("FileStoragePath", config.Values.FileStoragePath))
	} else {
		repo = &shorturlinmemory.InMemoryShortURLRepository{}
		logger.Log.Info("using in memory storage")
	}

	// ************************ HTTP ************************
	shortURLService := httpHandler.NewShortUrlsService(repo, auditPublisher)
	router := httpHandler.NewRouter(shortURLService)
	httpServer := &http.Server{
		Addr:    config.Values.ServerAddress,
		Handler: router,
	}
	if config.Values.EnableHTTPS {
		if config.Values.AutoCert {
			// конструируем менеджер TLS-сертификатов
			manager := &autocert.Manager{
				// директория для хранения сертификатов
				Cache: autocert.DirCache("shortener-cert-cache-dir"),
				// функция, принимающая Terms of Service издателя сертификатов
				Prompt: autocert.AcceptTOS,
			}

			httpServer.TLSConfig = manager.TLSConfig()
		}

		go func() {
			logger.Log.Info("starting tls server",
				zap.String("ServerAddress", config.Values.ServerAddress),
				zap.String("URLAddress", config.Values.URLAddress),
				zap.String("LogLevel", config.Values.LogLevel),
				zap.Bool("AutoCert", config.Values.AutoCert),
			)
			if err := httpServer.ListenAndServeTLS(config.Values.TLSCerfFile, config.Values.TLSKeyFile); err != nil && err != http.ErrServerClosed {
				log.Fatalf("failed listen and serve: %v", err)
			}
		}()
	} else {
		go func() {
			logger.Log.Info("starting server",
				zap.String("ServerAddress", config.Values.ServerAddress),
				zap.String("URLAddress", config.Values.URLAddress),
				zap.String("LogLevel", config.Values.LogLevel),
			)
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("failed listen and serve: %v", err)
			}
		}()
	}

	// ************************ gRPC ************************
	listen, err := net.Listen("tcp", config.Values.GRPCServerAddress)
	if err != nil {
		log.Fatalf("failed listen: %v", err)
	}
	tlsCreds, err := grpcHandler.LoadTLSCredentials()
	if err != nil {
		logger.Log.Debug("could not get TLS credentials", zap.Error(err))
	}
	var grpcServer *grpc.Server
	chainUnaryInterceptor := grpc.ChainUnaryInterceptor(
		interceptors.LoggingUnaryInterceptor,
		interceptors.AuthUnaryInterceptor,
	)
	if tlsCreds != nil {
		grpcServer = grpc.NewServer(
			grpc.Creds(tlsCreds),
			chainUnaryInterceptor,
		)
	} else {
		grpcServer = grpc.NewServer(
			chainUnaryInterceptor,
		)
	}

	pb.RegisterShortenerServiceServer(grpcServer, &grpcHandler.ShortenerServer{
		UnimplementedShortenerServiceServer: pb.UnimplementedShortenerServiceServer{},
		Repo:                                repo,
		AuditPublisher:                      auditPublisher,
	})
	go func() {
		logger.Log.Info("starting gRPC server",
			zap.String("GRPCServerAddress", config.Values.GRPCServerAddress),
		)
		if err := grpcServer.Serve(listen); err != nil && err != grpc.ErrServerStopped {
			log.Fatalf("failed to listen and serve grpc: %v", err)
		}
	}()

	// Wait for signal
	<-rootCtx.Done()
	stopRootCtx()
	logger.Log.Info("received shutdown signal, shutting down")

	shutdownCtx, cancelShutdownCtx := context.WithTimeout(context.Background(), shutdownPeriod)
	defer cancelShutdownCtx()

	grpcServerStop := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(grpcServerStop)
	}()

	select {
	case <-grpcServerStop:
		logger.Log.Info("gRPC server gracefully stopped")
	case <-shutdownCtx.Done():
		logger.Log.Debug("gRPC graceful shutdown timed out, forcing shutdown...")
		grpcServer.Stop()
	}

	err = httpServer.Shutdown(shutdownCtx)
	if err != nil {
		logger.Log.Warn("could not shutdown http server", zap.Error(err))
	}

	auditPublisher.Shutdown()

	logger.Log.Info("server shutdown complete")
}
