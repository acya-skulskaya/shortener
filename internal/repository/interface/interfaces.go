package interfaces

type ShortURLRepository interface {
	Get(id string) (originalURL string)
	Store(originalURL string) (id string)
}
