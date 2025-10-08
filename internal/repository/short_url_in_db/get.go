package shorturlindb

import (
	"context"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"go.uber.org/zap"
)

func (repo *InDBShortURLRepository) Get(ctx context.Context, id string) (originalURL string, err error) {
	row := repo.DB.QueryRowContext(ctx,
		"SELECT original_url, is_deleted FROM short_urls where id = $1", id)

	var isDeleted bool
	err = row.Scan(&originalURL, &isDeleted) // разбираем результат
	if err != nil {
		logger.Log.Debug("could not get row from db",
			zap.Error(err),
		)
		return "", err
	}

	if isDeleted {
		return "", errorsInternal.ErrIDDeleted
	}

	return originalURL, nil
}
