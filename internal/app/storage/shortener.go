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

func (container *LinksContainer) AddLink(f string, s string) {
	container.byShortMap[shortURL(s)] = fullURL(f)
	container.byFullMap[fullURL(f)] = shortURL(s)
}

func (container *LinksContainer) GetByShort(s string) (string, bool) {
	full, exist := container.byShortMap[shortURL(s)]
	return string(full), exist
}

func (container *LinksContainer) GetByFull(f string) (string, bool) {
	short, exist := container.byFullMap[fullURL(f)]
	return string(short), exist
}
