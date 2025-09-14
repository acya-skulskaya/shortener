package main

import (
	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/acya-skulskaya/shortener/internal/middleware"
	interfaces "github.com/acya-skulskaya/shortener/internal/repository/interface"
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

	shortURLService := NewShortUrlsService(&shorturljsonfile.JSONFileShortURLRepository{})

	router := NewRouter(shortURLService)
	err := http.ListenAndServe(config.Values.ServerAddress, router)

	logger.Log.Info("server started",
		zap.String("ServerAddress", config.Values.ServerAddress),
		zap.String("URLAddress", config.Values.URLAddress),
		zap.String("LogLevel", config.Values.LogLevel),
		zap.String("FileStoragePath", config.Values.FileStoragePath),
		zap.String("DatabaseDSN", config.Values.DatabaseDSN),
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
