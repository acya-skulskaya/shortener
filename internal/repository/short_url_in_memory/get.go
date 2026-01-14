package shorturlinmemory

import (
	"context"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
)

func (repo *InMemoryShortURLRepository) Get(ctx context.Context, id string) (originalURL string, err error) {
	item, err := cont.getItem(id)
	if err != nil {
		return "", err
	}

	if item.isDeleted == 1 {
		return "", errorsInternal.ErrIDDeleted
	}

	return item.originalURL, nil
}
