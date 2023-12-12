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
	loadErr := storage.LoadFromJSONFile(cfg.FileStoragePath)
	if loadErr != nil {
		log.Fatalf("Error encountered on restoring storage from file: %v", loadErr)
	}

	r := chi.NewRouter()
	registerMiddleware(r)
	registerRoutes(r, storage, cfg)

	err := http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		log.Fatalf("error ListenAndServe: %v", err)
	}
}

func registerRoutes(r *chi.Mux, storage storage.ShortenerStorage, cfg *config.Config) {
	postShortenHandler := handler.NewPostShortenHandler(storage, cfg)
	r.Handle("/", postShortenHandler)
	r.Method(http.MethodPost, "/api/shorten", postShortenHandler)
	r.Handle("/{id}", handler.NewGetURLHandler(storage, cfg))
}

func registerMiddleware(r *chi.Mux) {
	r.Use(middleware.ZapLogging)
	r.Use(middleware.GzipMiddleware)
}
