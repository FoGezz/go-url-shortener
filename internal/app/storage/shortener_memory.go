package storage

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

func (container *LinksMapping) AddLink(full string, short string) {
	container.byShortMap[shortURL(short)] = fullURL(full)
	container.byFullMap[fullURL(full)] = shortURL(short)
}

func (container *LinksMapping) GetByShort(s string) (full string, found bool) {
	f, exist := container.byShortMap[shortURL(s)]
	return string(f), exist
}

func (container *LinksMapping) GetByFull(f string) (short string, found bool) {
	s, exist := container.byFullMap[fullURL(f)]
	return string(s), exist
}
