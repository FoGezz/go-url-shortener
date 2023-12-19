package storage

import "context"

type fullURL string
type shortURL string

type shortToFullMap map[shortURL]fullURL
type fullToShortMap map[fullURL]shortURL

type LinksMapping struct {
	byShortMap shortToFullMap
	byFullMap  fullToShortMap
}

func NewLinksMapping() *LinksMapping {
	return &LinksMapping{
		byShortMap: make(shortToFullMap, 0),
		byFullMap:  make(fullToShortMap, 0),
	}
}

func (container *LinksMapping) AddLink(ctx context.Context, full string, short string) {
	container.byShortMap[shortURL(short)] = fullURL(full)
	container.byFullMap[fullURL(full)] = shortURL(short)
}

func (container *LinksMapping) GetByShort(ctx context.Context, s string) (full string, found bool) {
	f, exist := container.byShortMap[shortURL(s)]
	return string(f), exist
}

func (container *LinksMapping) GetByFull(ctx context.Context, f string) (short string, found bool) {
	s, exist := container.byFullMap[fullURL(f)]
	return string(s), exist
}
