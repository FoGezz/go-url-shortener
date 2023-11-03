package main

import (
	"log"
	"net/http"
	"os"

	"github.com/FoGezz/go-url-shortener/internal/app/handler"
	"github.com/FoGezz/go-url-shortener/internal/app/storage"
	"github.com/gorilla/mux"
)

func main() {

	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	log.SetOutput(file)

	container := storage.NewLinksContainer()

	mux := mux.NewRouter()
	mux.Handle("/", handler.NewPostShortenHandler(container))
	mux.Handle("/{id}", handler.NewGetURLHandler(container))
	http.ListenAndServe(":8080", mux)
}
