package shorturlindb

import (
	"context"

	"github.com/acya-skulskaya/shortener/internal/logger"
	"go.uber.org/zap"
)

func (repo *InDBShortURLRepository) GetInternalStats(ctx context.Context) (urls int, users int, err error) {
	row := repo.DB.QueryRowContext(ctx,
		"SELECT count(id) FROM short_urls where is_deleted=false")

	var countURLS int
	err = row.Scan(&countURLS)
	if err != nil {
		logger.Log.Debug("could not get url count from db",
			zap.Error(err),
		)
		return 0, 0, err
	}

	row = repo.DB.QueryRowContext(ctx,
		"SELECT count(DISTINCT user_id) FROM short_urls where is_deleted=false")

	var countUsers int
	err = row.Scan(&countUsers)
	if err != nil {
		logger.Log.Debug("could not get users count from db",
			zap.Error(err),
		)
		return 0, 0, err
	}

	return countURLS, countUsers, nil
}
