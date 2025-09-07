package main

import (
	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/acya-skulskaya/shortener/internal/middleware"
	interfaces "github.com/acya-skulskaya/shortener/internal/repository/interface"
	"github.com/acya-skulskaya/shortener/internal/repository/short_url_json_file"
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

//type Container struct {
//	mu        sync.Mutex
//	shortUrls map[string]string
//}
//
//func (c *Container) add(id string, value string) {
//	c.mu.Lock()
//	defer c.mu.Unlock()
//	c.shortUrls[id] = value
//}
//func (c *Container) getURL(id string) string {
//	c.mu.Lock()
//	defer c.mu.Unlock()
//	return c.shortUrls[id]
//}

// var ShortUrls = make(map[string]string)
//var Cont = Container{shortUrls: make(map[string]string)}

func main() {
	config.Init()

	if err := logger.Init(config.Values.LogLevel); err != nil {
		panic(err)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestLogger)
	router.Use(middleware.RequestCompressor)

	shortURLService := NewShortUrlsService(&short_url_json_file.JSONFileShortURLRepository{})

	router.Post("/", shortURLService.apiPageMain)
	router.Get("/{id}", shortURLService.apiPageByID)
	router.Post("/api/shorten", shortURLService.apiShorten)
	err := http.ListenAndServe(config.Values.ServerAddress, router)

	logger.Log.Info("server started",
		zap.String("ServerAddress", config.Values.ServerAddress),
		zap.String("URLAddress", config.Values.URLAddress),
		zap.String("LogLevel", config.Values.LogLevel),
		zap.String("FileStoragePath", config.Values.FileStoragePath),
	)

	if err != nil {
		panic(err)
	}
}
