package shorturljsonfile

import (
	"context"

	"github.com/acya-skulskaya/shortener/internal/config"
	"github.com/acya-skulskaya/shortener/internal/logger"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
	"go.uber.org/zap"
)

func (repo *JSONFileShortURLRepository) GetUserUrls(ctx context.Context, userID string) (list []jsonModel.BatchURLList, err error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	reader, err := NewFileReader(repo.FileStoragePath)
	if err != nil {
		logger.Log.Debug("could not create reader",
			zap.Error(err),
			zap.String("file", repo.FileStoragePath),
		)
		return nil, err
	}
	defer reader.Close()
	rows, err := reader.ReadFile()

	for _, row := range rows {
		if row.UserID == userID {
			if row.IsDeleted == 1 {
				continue
			}

			listItem := jsonModel.BatchURLList{
				OriginalURL: row.OriginalURL,
				ShortURL:    config.Values.URLAddress + "/" + row.ID,
			}

			list = append(list, listItem)
		}
	}

	return list, err
}
