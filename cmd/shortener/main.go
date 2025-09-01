package main

import (
	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"github.com/acya-skulskaya/shortener/internal/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type Container struct {
	mu        sync.Mutex
	shortUrls map[string]string
}

func (c *Container) add(id string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.shortUrls[id] = value
}
func (c *Container) getUrl(id string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.shortUrls[id]
}

// Cont TODO save urls in db
// TODO сделать репозиторий когда хранение будет в бд
// var ShortUrls = make(map[string]string)
var Cont = Container{shortUrls: make(map[string]string)}

func main() {
	config.Init()

	if err := logger.Init(config.Values.LogLevel); err != nil {
		panic(err)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestLogger)

	router.Post("/", apiPageMain)
	router.Get("/{id}", apiPageByID)
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
