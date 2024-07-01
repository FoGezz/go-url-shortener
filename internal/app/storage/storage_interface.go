package storage

import "context"

type ShortenerStorage interface {
	AddLink(ctx context.Context, full string, short string)
	GetByShort(ctx context.Context, s string) (full string, found bool, deleted bool)
	GetByFull(ctx context.Context, f string) (short string, found bool)
	LoadFromJSONFile(path string) error
	SaveJSONToFile(path string)
	GetByUserUUID(ctx context.Context, userUUID string) (*shortToFullMap, error)
	DeleteAsync(ctx context.Context, shortURLs []string, userUUID string)
}
