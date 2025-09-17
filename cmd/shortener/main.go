package main

import (
	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/acya-skulskaya/shortener/internal/middleware"
	interfaces "github.com/acya-skulskaya/shortener/internal/repository/interface"
	shorturlindb "github.com/acya-skulskaya/shortener/internal/repository/short_url_in_db"
	shorturlinmemory "github.com/acya-skulskaya/shortener/internal/repository/short_url_in_memory"
	shorturljsonfile "github.com/acya-skulskaya/shortener/internal/repository/short_url_json_file"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

type ShortUrlsService struct {
	repo interfaces.ShortURLRepository
}

func NewShortUrlsService(su interfaces.ShortURLRepository) *ShortUrlsService {
	return &ShortUrlsService{
		repo: su,
	}
}

func main() {
	config.Init()

	if err := logger.Init(config.Values.LogLevel); err != nil {
		panic(err)
	}

	shortURLService := NewShortUrlsService(&shorturlinmemory.InMemoryShortURLRepository{})
	if len(config.Values.DatabaseDSN) != 0 {
		db, err := shorturlindb.NewInDBShortURLRepository(config.Values.DatabaseDSN)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		shortURLService = NewShortUrlsService(&shorturlindb.InDBShortURLRepository{DB: db})
		logger.Log.Info("using db storage", zap.String("DatabaseDSN", config.Values.DatabaseDSN))
	} else if len(config.Values.FileStoragePath) != 0 {
		shortURLService = NewShortUrlsService(&shorturljsonfile.JSONFileShortURLRepository{})
		logger.Log.Info("using file storage", zap.String("FileStoragePath", config.Values.FileStoragePath))
	} else {
		logger.Log.Info("using in memory storage")
	}

	router := NewRouter(shortURLService)
	err := http.ListenAndServe(config.Values.ServerAddress, router)

	logger.Log.Info("server started",
		zap.String("ServerAddress", config.Values.ServerAddress),
		zap.String("URLAddress", config.Values.URLAddress),
		zap.String("LogLevel", config.Values.LogLevel),
	)

	if err != nil {
		panic(err)
	}
}

func NewRouter(su *ShortUrlsService) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestLogger)
	router.Use(middleware.RequestCompressor)

	router.Post("/", su.apiPageMain)
	router.Get("/{id}", su.apiPageByID)
	router.Get("/ping", su.apiPingDB)
	router.Post("/api/shorten", su.apiShorten)

	return router
}
