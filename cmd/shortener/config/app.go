package config

import (
	"github.com/FoGezz/go-url-shortener/internal/app/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Cfg     *Config
	Storage storage.ShortenerStorage
	DBPool  *pgxpool.Pool
}

func NewApp(cfg *Config, storage storage.ShortenerStorage) *App {
	return &App{cfg, storage, nil}
}
