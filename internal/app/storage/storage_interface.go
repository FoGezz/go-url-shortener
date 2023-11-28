package storage

type ShortenerStorage interface {
	AddLink(full string, short string)
	GetByShort(s string) (full string, found bool)
	GetByFull(f string) (short string, found bool)
	LoadFromJsonFile(path string)
	SaveJsonToFile(path string)
}
