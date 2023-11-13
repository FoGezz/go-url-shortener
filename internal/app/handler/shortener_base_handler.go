package handler

import (
	"github.com/FoGezz/go-url-shortener/cmd/shortener/config"
	"github.com/FoGezz/go-url-shortener/internal/app/storage"
)

type ShortenerHandler struct {
	storage storage.ShortenerStorage
	cfg     *config.Config
}
