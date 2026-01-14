package shorturlinmemory

import (
	"context"
	"errors"
	"fmt"
	"github.com/acya-skulskaya/shortener/internal/config"
	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
	"github.com/acya-skulskaya/shortener/internal/helpers"
	"github.com/acya-skulskaya/shortener/internal/logger"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
	"go.uber.org/zap"
	"sync"
)

type InMemoryShortURLRepository struct {
}

type Container struct {
	mu        sync.Mutex
	shortUrls map[string]shortURL
}

type shortURL struct {
	shortURL    string
	originalURL string
	userID      string
}

func (c *Container) add(id string, value string, userID string) (idAdded string, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.shortUrls[id]; ok {
		return id, errorsInternal.ErrConflictID
	}
	for i, item := range c.shortUrls {
		if item.originalURL == value {
			return i, errorsInternal.ErrConflictOriginalURL
		}
	}
	c.shortUrls[id] = shortURL{originalURL: value, userID: userID, shortURL: id}

	return id, nil
}

func (c *Container) getURL(id string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.shortUrls[id].originalURL
}

func (c *Container) getByUserID(userID string) (list []shortURL) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, item := range c.shortUrls {

		if item.userID == userID {
			list = append(list, item)
		}
	}

	return list
}

var cont = Container{shortUrls: make(map[string]shortURL)}

func (repo *InMemoryShortURLRepository) Get(ctx context.Context, id string) (originalURL string) {
	return cont.getURL(id)
}

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

func (repo *InMemoryShortURLRepository) Store(ctx context.Context, originalURL string, userID string) (id string, err error) {
	id = helpers.RandStringRunes(10)

	id, err = cont.add(id, originalURL, userID)
	if err != nil {
		if errors.Is(err, errorsInternal.ErrConflictOriginalURL) {
			return id, err
		} else {
			return "", err
		}
	}

	return id, nil
}

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
