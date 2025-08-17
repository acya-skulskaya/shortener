package main

import (
	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/go-chi/chi"
	"net/http"
)

// ShortUrls TODO save urls in db
var ShortUrls = make(map[string]string)

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
