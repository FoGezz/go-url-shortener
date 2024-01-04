package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStorage struct {
	conn *pgxpool.Conn
}

func NewDBStorage(c *pgxpool.Conn) *DBStorage {
	return &DBStorage{conn: c}
}

func (st *DBStorage) AddLink(ctx context.Context, full string, short string) {
	_, _ = st.conn.Exec(ctx, "INSERT INTO links(long,short) VALUES ($1,$2);", full, short)
}
func (st *DBStorage) GetByShort(ctx context.Context, s string) (full string, found bool) {
	err := st.conn.QueryRow(ctx, "SELECT long FROM links WHERE short = $1;", s).Scan(&full)
	if err != nil {
		return "", false
	}
	if full != "" {
		found = true
	}
	return
}

func (st *DBStorage) GetByFull(ctx context.Context, f string) (short string, found bool) {
	err := st.conn.QueryRow(ctx, "SELECT short FROM links WHERE long = $1;", f).Scan(&short)
	if err != nil {
		return "", false
	}
	if short != "" {
		found = true
	}
	return
}

func (st *DBStorage) LoadFromJSONFile(path string) error {
	return nil
}
func (st *DBStorage) SaveJSONToFile(path string) {
}
