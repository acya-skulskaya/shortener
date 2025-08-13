package main

import (
	"net/http"
)

var Urls map[string]string

func main() {
	Urls = make(map[string]string)

	mux := http.NewServeMux()
	mux.HandleFunc(`/{id}`, apiPageByID)
	mux.HandleFunc(`/`, apiPageMain)

	err := http.ListenAndServe(`:80`, mux)

	if err != nil {
		panic(err)
	}
}
