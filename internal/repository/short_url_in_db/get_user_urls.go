package shorturlindb

import (
	"context"

	"github.com/acya-skulskaya/shortener/internal/logger"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
	"go.uber.org/zap"
)

func (repo *InDBShortURLRepository) GetUserUrls(ctx context.Context, userID string) (list []jsonModel.BatchURLList, err error) {
	rows, err := repo.DB.QueryContext(ctx, "SELECT short_url, original_url from short_urls WHERE user_id = $1 AND is_deleted = false", userID)
	if err != nil {
		logger.Log.Debug("could not query from db", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		listShortenItem := jsonModel.BatchURLList{}
		err = rows.Scan(&listShortenItem.ShortURL, &listShortenItem.OriginalURL)
		if err != nil {
			logger.Log.Debug("could not scan row", zap.Error(err))
			return nil, err
		}

		list = append(list, listShortenItem)
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Debug("error getting user urls", zap.Error(err))
		return nil, err
	}

	return list, err
}
