package main

import (
	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/go-chi/chi"
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

// Cont TODO save urls in db
// TODO сделать репозиторий когда хранение будет в бд
// var ShortUrls = make(map[string]string)
var Cont = Container{shortUrls: make(map[string]string)}

func main() {
	config.Init()

	router := chi.NewRouter()
	router.Post("/", apiPageMain)
	router.Get("/{id}", apiPageByID)
	err := http.ListenAndServe(config.Values.ServerAddress, router)

	if err != nil {
		panic(err)
	}
}
