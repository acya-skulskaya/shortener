package shorturlinmemory

import (
	"context"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"go.uber.org/zap"
)

func (repo *InMemoryShortURLRepository) DeleteUserUrls(ctx context.Context, list []string, userID string) (err error) {
	for _, id := range list {
		item, err := cont.getItem(id)
		if err != nil {
			return err
		}
		if item.userID != userID {
			return errorsInternal.ErrUserIDUnauthorized
		}
		if item.isDeleted == 1 {
			return errorsInternal.ErrIDDeleted
		}
	}

	for _, id := range list {
		go func(id string) {
			err = cont.deleteItem(id)
			if err != nil {
				logger.Log.Debug("could not delete item", zap.String("id", id), zap.Error(err))
			} else {
				logger.Log.Info("item was deleted", zap.String("id", id))
			}
		}(id)
	}

	return nil
}
