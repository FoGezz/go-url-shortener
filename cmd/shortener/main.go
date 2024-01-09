package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/FoGezz/go-url-shortener/cmd/shortener/config"
	"github.com/FoGezz/go-url-shortener/internal/app/handler"
	"github.com/FoGezz/go-url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := &config.Config{}
	cfg.Load()
	cfg.DBDSN = "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable"
	fmt.Println("Running on conf", cfg)

	db := storage.NewLinksMapping()
	loadErr := db.LoadFromJSONFile(cfg.FileStoragePath)
	if loadErr != nil {
		log.Fatalf("Error encountered on restoring storage from file: %v", loadErr)
	}
	app := config.NewApp(cfg, db)
	if cfg.DBDSN != "" {
		pool, DBErr := storage.NewDB(cfg.DBDSN)
		if DBErr != nil {
			log.Fatalf("Error encountered on connecting to DB: %v", DBErr)
		}
		app.DBPool = pool
		conn, DBErr := pool.Acquire(context.Background())
		if DBErr != nil {
			log.Fatalf("Error encountered on acquiring conn: %v", DBErr)
		}
		app.Storage = storage.NewDBStorage(conn)
	}

	r := chi.NewRouter()
	handler.RegisterMiddleware(r)
	handler.RegisterRoutes(r, app)

	err := http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		log.Fatalf("error ListenAndServe: %v", err)
	}
}
