package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/FoGezz/go-url-shortener/cmd/shortener/config"
	"github.com/FoGezz/go-url-shortener/internal/app/handler"
	"github.com/FoGezz/go-url-shortener/internal/app/middleware"
	"github.com/FoGezz/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := &config.Config{}
	cfg.Load()
	fmt.Println("Running on conf", cfg)

	storage := storage.NewLinksMapping()
	storage.LoadFromJsonFile(cfg.FileStoragePath)

	r := chi.NewRouter()
	r.Use(middleware.ZapLogging)
	r.Use(middleware.GzipMiddleware)
	postShortenHandler := handler.NewPostShortenHandler(storage, cfg)
	r.Handle("/", postShortenHandler)
	r.Method(http.MethodPost, "/api/shorten", postShortenHandler)
	r.Handle("/{id}", handler.NewGetURLHandler(storage, cfg))

	err := http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		log.Fatalf("error ListenAndServe: %v", err)
	}
}
