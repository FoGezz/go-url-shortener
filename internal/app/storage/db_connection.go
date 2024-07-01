package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/FoGezz/go-url-shortener/internal/app/storage/migrations"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func runMigrations(DSN string) error {
	d, err := iofs.New(migrations.MigrationsDir, ".")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, DSN)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	return nil
}

func NewDB(DSN string) (*pgxpool.Pool, error) {
	err := runMigrations(DSN)
	if err != nil {
		return nil, fmt.Errorf("error on runMigrations: %w", err)
	}
	if DSN != "" {
		conn, err := pgxpool.New(context.Background(), DSN)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}

	return nil, nil
}
