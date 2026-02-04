package shorturljsonfile

import (
	"context"
	"slices"

	"github.com/acya-skulskaya/shortener/internal/logger"
	"go.uber.org/zap"
)

func (repo *JSONFileShortURLRepository) DeleteUserUrls(ctx context.Context, list []string, userID string) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	reader := NewFileReader(repo.FileStoragePath)
	existingRows, err := reader.ReadFile()
	if err != nil {
		logger.Log.Debug("could not read file",
			zap.Error(err),
			zap.String("file", repo.FileStoragePath),
		)
		return
	}

	writer := NewFileWriter(repo.FileStoragePath)

	doOverwrite := false
	for i, existingRow := range existingRows {
		if !slices.Contains(list, existingRow.ID) {
			continue
		}

		if existingRow.UserID != userID {
			logger.Log.Debug("id belongs to another user", zap.String("id", existingRow.ID), zap.String("userID", userID))
			continue
		}
		if existingRow.IsDeleted == 1 {
			logger.Log.Debug("id is already deleted", zap.String("id", existingRow.ID), zap.String("userID", userID))
			continue
		}

		doOverwrite = true

		existingRow.IsDeleted = 1
		existingRows[i] = existingRow
		logger.Log.Info("item was marked as deleted", zap.String("id", existingRow.ID))
	}

	if doOverwrite {
		err = writer.OverwriteFile(existingRows)
		if err != nil {
			logger.Log.Debug("could not delete ulrs in file",
				zap.Error(err),
				zap.Any("list", list),
			)
			return
		}
		logger.Log.Info("items were deleted", zap.Any("list", list))
	}
}
