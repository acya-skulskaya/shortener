package shorturljsonfile

import (
	"context"

	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	"go.uber.org/zap"
)

func (repo *JSONFileShortURLRepository) Get(ctx context.Context, id string) (originalURL string, err error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	reader, err := NewFileReader(repo.FileStoragePath)
	if err != nil {
		logger.Log.Debug("could not create reader",
			zap.Error(err),
			zap.String("id", id),
			zap.String("file", repo.FileStoragePath),
		)
		return "", err
	}
	defer reader.Close()

	list, err := reader.ReadFile()
	if err != nil {
		logger.Log.Debug("could not read file",
			zap.Error(err),
			zap.String("id", id),
			zap.String("file", repo.FileStoragePath),
		)
		return "", err
	}

	for _, l := range list {
		if l.ID == id {
			if l.IsDeleted == 1 {
				return "", errorsInternal.ErrIDDeleted
			}
			return l.OriginalURL, nil
		}
	}

	return "", nil
}
