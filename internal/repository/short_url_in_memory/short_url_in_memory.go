package shorturlinmemory

import (
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
	shortUrls map[string]string
}

func (c *Container) add(id string, value string) (idAdded string, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.shortUrls[id]; ok {
		return id, errorsInternal.ErrConflictID
	}
	for i, item := range c.shortUrls {
		if item == value {
			return i, errorsInternal.ErrConflictOriginalURL
		}
	}
	c.shortUrls[id] = value

	return id, nil
}

func (c *Container) getURL(id string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.shortUrls[id]
}

var cont = Container{shortUrls: make(map[string]string)}

func (repo *InMemoryShortURLRepository) Get(id string) (originalURL string) {
	return cont.getURL(id)
}

func (repo *InMemoryShortURLRepository) Store(originalURL string) (id string, err error) {
	id = helpers.RandStringRunes(10)

	id, err = cont.add(id, originalURL)
	if err != nil {
		if err != nil && errors.Is(err, errorsInternal.ErrConflictOriginalURL) {
			return id, err
		} else {
			return "", err
		}
	}

	return id, nil
}

func (repo *InMemoryShortURLRepository) StoreBatch(listOriginal []jsonModel.BatchURLList) (listShorten []jsonModel.BatchURLList, err error) {
	for _, item := range listOriginal {
		id, err := cont.add(item.CorrelationID, item.OriginalURL)
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
