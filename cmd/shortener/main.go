package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/FoGezz/go-url-shortener/cmd/shortener/config"
	"github.com/FoGezz/go-url-shortener/internal/app/handler"
	"github.com/FoGezz/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := config.Config{}
	cfg.ParseFlags()
	fmt.Println("Run on conf", cfg)

	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	log.SetOutput(file)

	storage := storage.NewLinksContainer()
	r := chi.NewRouter()
	r.Handle("/", handler.NewPostShortenHandler(storage, cfg))
	r.Handle("/{id}", handler.NewGetURLHandler(storage, cfg))

	http.ListenAndServe(cfg.ServerAddress, r)
}
