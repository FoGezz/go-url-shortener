package storage

type fullURL string
type shortURL string

type shortToFullMap map[shortURL]fullURL
type fullToShortMap map[fullURL]shortURL

type LinksContainer struct {
	byShortMap shortToFullMap
	byFullMap  fullToShortMap
}

func NewLinksContainer() *LinksContainer {
	return &LinksContainer{
		byShortMap: make(shortToFullMap, 0),
		byFullMap:  make(fullToShortMap, 0),
	}
}

func (container *LinksContainer) AddLink(full string, short string) {
	container.byShortMap[shortURL(short)] = fullURL(full)
	container.byFullMap[fullURL(full)] = shortURL(short)
}

func (container *LinksContainer) GetByShort(s string) (full string, found bool) {
	f, exist := container.byShortMap[shortURL(s)]
	return string(f), exist
}

func (container *LinksContainer) GetByFull(f string) (short string, found bool) {
	s, exist := container.byFullMap[fullURL(f)]
	return string(s), exist
}
