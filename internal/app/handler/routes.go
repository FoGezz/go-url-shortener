package handler

import (
	"net/http"

	"github.com/FoGezz/go-url-shortener/cmd/shortener/config"
	"github.com/FoGezz/go-url-shortener/internal/app/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r *chi.Mux, app *config.App) {
	postShortenHandler := NewPostShortenHandler(app)
	r.Handle("/", postShortenHandler)
	r.Method(http.MethodPost, "/api/shorten", postShortenHandler)
	r.Handle("/{id}", NewGetURLHandler(app))
	r.Method(http.MethodGet, "/ping", NewGetPingHandler(app))
	r.Method(http.MethodPost, "/api/shorten/batch", NewPostShortenBatchHandler(app))
}

func RegisterMiddleware(r *chi.Mux) {
	r.Use(middleware.ZapLogging)
	r.Use(middleware.GzipMiddleware)
}
