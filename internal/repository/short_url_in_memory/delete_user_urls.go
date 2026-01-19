package shorturlinmemory

import (
	"context"

	"github.com/acya-skulskaya/shortener/internal/logger"
	"go.uber.org/zap"
)

func (repo *InMemoryShortURLRepository) DeleteUserUrls(ctx context.Context, list []string, userID string) {
	for _, id := range list {
		go func(id string, userID string) {
			item, errGet := cont.getItem(id)
			if errGet != nil {
				logger.Log.Debug("could not get item", zap.String("id", id), zap.String("userID", userID), zap.Error(errGet))
				return
			}
			if item.userID != userID {
				logger.Log.Debug("id belongs to anther user", zap.String("id", id), zap.String("userID", userID))
				return
			}
			if item.isDeleted == 1 {
				logger.Log.Debug("id os already deleted", zap.String("id", id), zap.String("userID", userID))
				return
			}

			errDel := cont.deleteItem(id)
			if errDel != nil {
				logger.Log.Debug("could not delete item", zap.String("id", id), zap.String("userID", userID), zap.Error(errDel))
			} else {
				logger.Log.Info("item was deleted", zap.String("id", id), zap.String("userID", userID))
			}
		}(id, userID)
	}
}
