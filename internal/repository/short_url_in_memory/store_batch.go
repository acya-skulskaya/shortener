package shorturlinmemory

import (
	"context"
	"errors"
	"fmt"
	"github.com/acya-skulskaya/shortener/internal/config"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/logger"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
	"go.uber.org/zap"
)

func (repo *InMemoryShortURLRepository) StoreBatch(ctx context.Context, listOriginal []jsonModel.BatchURLList, userID string) (listShorten []jsonModel.BatchURLList, err error) {
	for _, item := range listOriginal {
		id, err := cont.add(item.CorrelationID, item.OriginalURL, userID)
		listShortenItem := jsonModel.BatchURLList{
			CorrelationID: item.CorrelationID,
			ShortURL:      config.Values.URLAddress + "/" + item.CorrelationID,
		}

		if err != nil {
			if errors.Is(err, errorsInternal.ErrConflictOriginalURL) || errors.Is(err, errorsInternal.ErrConflictID) {
				listShortenItem.Err = fmt.Sprint(err)
				if errors.Is(err, errorsInternal.ErrConflictID) {
					listShortenItem.CorrelationID = id
				}
			} else {
				logger.Log.Debug("could not add item",
					zap.Error(err),
					zap.Any("item", item),
				)
				return nil, err
			}
		}

		listShorten = append(listShorten, listShortenItem)
	}

	return listShorten, err
}
