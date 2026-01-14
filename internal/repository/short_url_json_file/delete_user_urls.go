package shorturljsonfile

import (
	"context"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
	"go.uber.org/zap"
	"slices"
)

func checkAndGetIdsInFile(repo *JSONFileShortURLRepository, list []string, userID string) (existingRows []jsonModel.URLList, err error) {

	for _, id := range list {
		idExists := false
		for _, existingRow := range existingRows {
			if existingRow.ID != id {
				continue
			}
			idExists = true
			if existingRow.UserID != userID {
				return []jsonModel.URLList{}, errorsInternal.ErrUserIDUnauthorized
			}
			if existingRow.IsDeleted == 1 {
				return []jsonModel.URLList{}, errorsInternal.ErrIDDeleted
			}
		}

		if !idExists {
			return []jsonModel.URLList{}, errorsInternal.ErrIDNotFound
		}
	}

	return existingRows, nil
}

func (repo *JSONFileShortURLRepository) DeleteUserUrls(ctx context.Context, list []string, userID string) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	reader, err := NewFileReader(repo.FileStoragePath)
	if err != nil {
		logger.Log.Debug("could not create reader",
			zap.Error(err),
			zap.String("file", repo.FileStoragePath),
		)
		return
	}
	defer reader.Close()
	existingRows, err := reader.ReadFile()
	if err != nil {
		logger.Log.Debug("could not read file",
			zap.Error(err),
			zap.String("file", repo.FileStoragePath),
		)
		return
	}

	writer, err := NewFileWriter(repo.FileStoragePath)
	if err != nil {
		logger.Log.Debug("could not create file writer",
			zap.Error(err),
		)
		return
	}

	doOverwrite := false
	for i, existingRow := range existingRows {
		if !slices.Contains(list, existingRow.ID) {
			continue
		}

		if existingRow.UserID != userID {
			logger.Log.Debug("id belongs to anther user", zap.String("id", existingRow.ID), zap.String("userID", userID))
			continue
		}
		if existingRow.IsDeleted == 1 {
			logger.Log.Debug("id os already deleted", zap.String("id", existingRow.ID), zap.String("userID", userID))
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
