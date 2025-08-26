package main

import (
	"github.com/go-chi/chi"
	"net/http"
)

var ShortUrls map[string]string

func main() {
	ShortUrls = make(map[string]string)

	router := chi.NewRouter()

	router.Post("/", apiPageMain)
	router.Get("/{id}", apiPageByID)

	err := http.ListenAndServe(":8080", router)

	if err != nil {
		panic(err)
	}
}
