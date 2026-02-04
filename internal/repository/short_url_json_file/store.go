package shorturljsonfile

import (
	"context"
	"errors"

	"github.com/acya-skulskaya/shortener/internal/config"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/helpers"
	"github.com/acya-skulskaya/shortener/internal/logger"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
	"go.uber.org/zap"
)

func (repo *JSONFileShortURLRepository) Store(ctx context.Context, originalURL string, userID string) (id string, err error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	reader := NewFileReader(repo.FileStoragePath)
	existingRows, err := reader.ReadFile()
	if err != nil {
		logger.Log.Debug("could not read file",
			zap.Error(err),
			zap.String("file", repo.FileStoragePath),
		)
		return "", err
	}

	for _, existingRow := range existingRows {
		if existingRow.OriginalURL == originalURL {
			return existingRow.ID, errorsInternal.ErrConflictOriginalURL
		}
	}

	id = helpers.RandStringRunes(10)

	row := jsonModel.URLList{
		ID:          id,
		ShortURL:    config.Values.URLAddress + "/" + id,
		OriginalURL: originalURL,
		UserID:      userID,
	}

	writer := NewFileWriter(repo.FileStoragePath)
	err = writer.WriteFile(row)
	if err != nil {
		if errors.Is(err, errorsInternal.ErrConflictOriginalURL) || errors.Is(err, errorsInternal.ErrConflictID) {
			return id, err
		}

		logger.Log.Debug("could not write short url to file",
			zap.Error(err),
			zap.String("id", id),
			zap.String("url", originalURL),
			zap.String("file", repo.FileStoragePath),
		)
		return "", err
	}

	return id, nil
}
