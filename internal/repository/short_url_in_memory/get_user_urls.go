package shorturlinmemory

import (
	"context"

	"github.com/acya-skulskaya/shortener/internal/config"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
)

func (repo *InMemoryShortURLRepository) GetUserUrls(ctx context.Context, userID string) (list []jsonModel.BatchURLList, err error) {
	shortURLs := cont.getByUserID(userID)

	for _, item := range shortURLs {
		listItem := jsonModel.BatchURLList{
			OriginalURL: item.originalURL,
			ShortURL:    config.Values.URLAddress + "/" + item.shortURL,
		}

		list = append(list, listItem)
	}

	return list, nil
}
