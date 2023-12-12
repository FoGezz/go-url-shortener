package main

import (
	"net/http"

	"github.com/FoGezz/go-url-shortener/cmd/shortener/config"
	"github.com/FoGezz/go-url-shortener/internal/app/handler"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(r *chi.Mux, app *config.App) {
	postShortenHandler := handler.NewPostShortenHandler(app)
	r.Handle("/", postShortenHandler)
	r.Method(http.MethodPost, "/api/shorten", postShortenHandler)
	r.Handle("/{id}", handler.NewGetURLHandler(app))
	r.Method(http.MethodGet, "/ping", handler.NewGetPingHandler(app))
}
