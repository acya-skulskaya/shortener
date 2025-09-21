package interfaces

import jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"

type ShortURLRepository interface {
	Get(id string) (originalURL string)
	Store(originalURL string) (id string, err error)
	StoreBatch(listOriginal []jsonModel.BatchURLList) (listShorten []jsonModel.BatchURLList)
}
