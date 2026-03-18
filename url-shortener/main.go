package main

import (
	"fmt"
	"net/http"
	"url-shortener/handler"
	"url-shortener/storage"
)

func main() {
	s := storage.NewMemoryStorage()
	h := handler.NewHandler(s)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", h.Shorten)
	mux.HandleFunc("GET /{code}", h.Redirect)
	mux.HandleFunc("GET /stats/{code}", h.Stats)

	http.ListenAndServe(":8080", mux)
}

func printAll(s storage.Storage) {
	entries := s.All()

	for _, entry := range entries {
		fmt.Println(entry.String())
	}

}
