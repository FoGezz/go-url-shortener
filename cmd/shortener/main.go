package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/FoGezz/go-url-shortener/cmd/shortener/config"
	"github.com/FoGezz/go-url-shortener/internal/app/middleware"
	"github.com/FoGezz/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	app := config.NewApp(cfg, storage)

	if cfg.DbDSN != "" {
		conn, err := pgxpool.New(context.Background(), cfg.DbDSN)
		if err != nil {
			log.Fatalf("Error connecting to PG: %v", err)
		}
		app.DBPool = conn

		defer conn.Close()
	}

	r := chi.NewRouter()
	registerMiddleware(r)
	registerRoutes(r, app)

	err := http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		log.Fatalf("error ListenAndServe: %v", err)
	}
}

func registerMiddleware(r *chi.Mux) {
	r.Use(middleware.ZapLogging)
	r.Use(middleware.GzipMiddleware)
}
