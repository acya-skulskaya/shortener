package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/{id}/`, apiPageById)
	mux.HandleFunc(`/`, apiPageMain)

	err := http.ListenAndServe(`:80`, mux)

	if err != nil {
		panic(err)
	}
}
