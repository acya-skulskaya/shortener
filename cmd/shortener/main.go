package main

import (
	"net/http"
)

var ShortUrls map[string]string

func main() {
	ShortUrls = make(map[string]string)

	mux := http.NewServeMux()
	mux.HandleFunc(`/{id}`, apiPageByID)
	mux.HandleFunc(`/`, apiPageMain)

	err := http.ListenAndServe(`:8080`, mux)
	//err := http.ListenAndServe(`localhost`, mux)

	if err != nil {
		panic(err)
	}
}
