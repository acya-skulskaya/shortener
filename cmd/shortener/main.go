package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/acya-skulskaya/shortener/internal/middleware"
	"github.com/acya-skulskaya/shortener/internal/observer/audit/publisher"
	"github.com/acya-skulskaya/shortener/internal/observer/audit/subscribers"
	interfaces "github.com/acya-skulskaya/shortener/internal/repository/interface"
	shorturlindb "github.com/acya-skulskaya/shortener/internal/repository/short_url_in_db"
	shorturlinmemory "github.com/acya-skulskaya/shortener/internal/repository/short_url_in_memory"
	shorturljsonfile "github.com/acya-skulskaya/shortener/internal/repository/short_url_json_file"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
)

const (
	shutdownPeriod = 20 * time.Second
)

// ShortUrlsService provides access to URL Shortener storage interface and audit publisher
type ShortUrlsService struct {
	Repo           interfaces.ShortURLRepository
	auditPublisher publisher.Publisher
}

// NewShortUrlsService creates a new instance of ShortUrlsService
func NewShortUrlsService(su interfaces.ShortURLRepository, ap publisher.Publisher) *ShortUrlsService {
	return &ShortUrlsService{
		Repo:           su,
		auditPublisher: ap,
	}
}

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

	var shortURLService *ShortUrlsService
	if len(config.Values.DatabaseDSN) != 0 {
		db, err := shorturlindb.NewInDBShortURLRepository(config.Values.DatabaseDSN)
		if err != nil {
			logger.Log.Fatal("failed to init db storage", zap.Error(err))
		}
		defer db.Close()
		shortURLService = NewShortUrlsService(&shorturlindb.InDBShortURLRepository{DB: db}, auditPublisher)
		logger.Log.Info("using db storage", zap.String("DatabaseDSN", config.Values.DatabaseDSN))
	} else if len(config.Values.FileStoragePath) != 0 {
		shortURLService = NewShortUrlsService(&shorturljsonfile.JSONFileShortURLRepository{FileStoragePath: config.Values.FileStoragePath}, auditPublisher)
		logger.Log.Info("using file storage", zap.String("FileStoragePath", config.Values.FileStoragePath))
	} else {
		shortURLService = NewShortUrlsService(&shorturlinmemory.InMemoryShortURLRepository{}, auditPublisher)
		logger.Log.Info("using in memory storage")
	}

	router := NewRouter(shortURLService)

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

	// Wait for signal
	<-rootCtx.Done()
	stopRootCtx()
	logger.Log.Info("received shutdown signal, shutting down")

	shutdownCtx, cancelShutdownCtx := context.WithTimeout(context.Background(), shutdownPeriod)
	defer cancelShutdownCtx()
	err := httpServer.Shutdown(shutdownCtx)
	if err != nil {
		logger.Log.Warn("could not shutdown http server", zap.Error(err))
	}

	auditPublisher.Shutdown()

	logger.Log.Info("server shutdown complete")
}

// NewRouter initiates a new router with API's endpoints:
//   - POST / - creates a shortened URL
//   - GET /{id} — redirects to the original URL
//   - GET /ping — tests connection to DB
//   - POST /api/shorten - creates a shortened URL
//   - POST /api/shorten/batch - creates a  list of shortened URLs
//   - GET /api/user/urls - returns a list of URLs that were added by an authenticated user
//   - DELETE /api/user/urls - deletes a list of shortened URLs
func NewRouter(su *ShortUrlsService) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestLogger)
	router.Use(middleware.RequestCompressor)
	router.Use(middleware.CookieAuth)

	// pprof
	//router.Route("/debug/pprof", func(r chi.Router) {
	//	r.Handle("/", http.HandlerFunc(pprof.Index))
	//	r.Handle("/profile", http.HandlerFunc(pprof.Profile))
	//	r.Handle("/symbol", http.HandlerFunc(pprof.Symbol))
	//	r.Handle("/cmdline", http.HandlerFunc(pprof.Cmdline))
	//	r.Handle("/heap", pprof.Handler("heap"))
	//})

	router.Post("/", su.apiPageMain)
	router.Get("/{id}", su.apiPageByID)
	router.Get("/ping", su.apiPingDB)
	router.Post("/api/shorten", su.apiShorten)
	router.Post("/api/shorten/batch", su.apiShortenBatch)
	router.Get("/api/user/urls", su.apiUserURLs)
	router.Delete("/api/user/urls", su.apiDeleteUserURLs)
	router.Get("/api/internal/stats", su.apiInternalStats)

	return router
}
