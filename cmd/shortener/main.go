package main

import (
	"log"
	"net/http"
	"os"

	"github.com/FoGezz/go-url-shortener/internal/app/handler"
	"github.com/FoGezz/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

func main() {

	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	log.SetOutput(file)

	storage := storage.NewLinksContainer()
	r := chi.NewRouter()
	r.Handle("/", handler.NewPostShortenHandler(storage))
	r.Handle("/{id}", handler.NewGetURLHandler(storage))

	http.ListenAndServe(":8080", r)
}
