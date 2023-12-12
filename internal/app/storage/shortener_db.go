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
	r, err := st.conn.Query(context.Background(), "SELECT long,short FROM links WHERE short = $1;", s)
	if err != nil {
		return "", false
	}
	r.Scan(&full)
	if full != "" {
		found = true
	}
	return
}

func (st *DBStorage) GetByFull(f string) (short string, found bool) {
	r, err := st.conn.Query(context.Background(), "SELECT long,short FROM links WHERE long = $1;", f)
	if err != nil {
		return "", false
	}
	r.Scan(&short)
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
