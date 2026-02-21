package shorturljsonfile

import (
	"context"
	"slices"

	"github.com/acya-skulskaya/shortener/internal/logger"
	"go.uber.org/zap"
)

func (repo *JSONFileShortURLRepository) GetInternalStats(ctx context.Context) (urls int, users int, err error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	reader := NewFileReader(repo.FileStoragePath)

	list, err := reader.ReadFile()
	if err != nil {
		logger.Log.Debug("could not read file",
			zap.Error(err),
			zap.String("file", repo.FileStoragePath),
		)
		return 0, 0, err
	}

	var countURLs int
	var uniqueUsers []string

	for _, l := range list {
		if l.IsDeleted == 1 {
			continue
		}
		countURLs++

		if !slices.Contains(uniqueUsers, l.UserID) {
			uniqueUsers = append(uniqueUsers, l.UserID)
		}
	}

	return countURLs, len(uniqueUsers), nil
}
