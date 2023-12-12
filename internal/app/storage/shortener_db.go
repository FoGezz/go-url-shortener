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

func (st *DBStorage) AddLink(full string, short string) {
	_, _ = st.conn.Exec(context.Background(), "INSERT INTO links(long,short) VALUES ($1,$2);", full, short)
}
func (st *DBStorage) GetByShort(s string) (full string, found bool) {
	err := st.conn.QueryRow(context.Background(), "SELECT long FROM links WHERE short = $1;", s).Scan(&full)
	if err != nil {
		return "", false
	}
	if full != "" {
		found = true
	}
	return
}

func (st *DBStorage) GetByFull(f string) (short string, found bool) {
	err := st.conn.QueryRow(context.Background(), "SELECT short FROM links WHERE long = $1;", f).Scan(&short)
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
