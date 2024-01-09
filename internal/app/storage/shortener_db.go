package storage

import (
	"context"
	"errors"
	"log"

	"github.com/FoGezz/go-url-shortener/internal/app/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStorage struct {
	conn *pgxpool.Conn
}

func NewDBStorage(c *pgxpool.Conn) *DBStorage {
	return &DBStorage{conn: c}
}

func (st *DBStorage) AddLink(ctx context.Context, full string, short string) {
	_, err := st.conn.Exec(ctx, "INSERT INTO links(long,short,user_uuid) VALUES ($1,$2,$3);", full, short, ctx.Value(middleware.USER_ID_KEY))
	if err != nil {
		log.Println(err)
	}
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

func (st *DBStorage) GetByUserUUID(ctx context.Context, userUUID string) (*shortToFullMap, error) {
	rows, err := st.conn.Query(ctx, "SELECT short,long FROM links WHERE user_uuid = $1;", userUUID)
	if err != nil {
		return nil, err
	}
	m := shortToFullMap{}
	for rows.Next() {
		var short, long string
		err := rows.Scan(&short, &long)
		if err != nil {
			return nil, errors.Join(errors.New("error scanning rows from rowset"), err)
		}
		m[shortURL(short)] = fullURL(long)
	}

	return &m, nil
}

func (st *DBStorage) LoadFromJSONFile(path string) error {
	return nil
}
func (st *DBStorage) SaveJSONToFile(path string) {
}
