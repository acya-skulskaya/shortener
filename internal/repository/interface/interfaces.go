package interfaces

import (
	"context"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
)

type ShortURLRepository interface {
	Get(ctx context.Context, id string) (originalURL string)
	Store(ctx context.Context, originalURL string) (id string, err error)
	StoreBatch(ctx context.Context, listOriginal []jsonModel.BatchURLList) (listShorten []jsonModel.BatchURLList, err error)
}
