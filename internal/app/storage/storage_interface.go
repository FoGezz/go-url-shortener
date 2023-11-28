package storage

type ShortenerStorage interface {
	AddLink(full string, short string)
	GetByShort(s string) (full string, found bool)
	GetByFull(f string) (short string, found bool)
	LoadFromJSONFile(path string)
	SaveJSONToFile(path string)
}
