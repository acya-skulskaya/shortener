package interfaces

import (
	"context"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
)

type ShortURLRepository interface {
	Get(ctx context.Context, id string) (originalURL string, err error)
	GetUserUrls(ctx context.Context, userID string) (list []jsonModel.BatchURLList, err error)
	DeleteUserUrls(ctx context.Context, list []string, userID string) (err error)
	Store(ctx context.Context, originalURL string, userID string) (id string, err error)
	StoreBatch(ctx context.Context, listOriginal []jsonModel.BatchURLList, userID string) (listShorten []jsonModel.BatchURLList, err error)
}
