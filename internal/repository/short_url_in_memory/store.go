package shorturlinmemory

import (
	"context"
	"errors"

	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/helpers"
)

func (repo *InMemoryShortURLRepository) Store(ctx context.Context, originalURL string, userID string) (id string, err error) {
	id = helpers.RandStringRunes(10)

	id, err = cont.add(id, originalURL, userID)
	if err != nil {
		if errors.Is(err, errorsInternal.ErrConflictOriginalURL) {
			return id, err
		} else {
			return "", err
		}
	}

	return id, nil
}
