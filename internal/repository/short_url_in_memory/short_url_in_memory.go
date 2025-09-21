package shorturlinmemory

import (
	"fmt"
	"github.com/acya-skulskaya/shortener/internal/config"
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

func (c *Container) add(id string, value string) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.shortUrls[id]; ok {
		return fmt.Errorf("short url with id %s already exists", id)
	}
	c.shortUrls[id] = value

	return nil
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

	err = cont.add(id, originalURL)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (repo *InMemoryShortURLRepository) StoreBatch(listOriginal []jsonModel.BatchURLList) (listShorten []jsonModel.BatchURLList) {
	for _, item := range listOriginal {
		err := cont.add(item.CorrelationID, item.OriginalURL)
		if err != nil {
			logger.Log.Debug("could not add item",
				zap.Error(err),
				zap.Any("item", item),
			)
			return nil
		}

		listShorten = append(listShorten, jsonModel.BatchURLList{
			CorrelationID: item.CorrelationID,
			ShortURL:      config.Values.URLAddress + "/" + item.CorrelationID,
		})
	}

	return listShorten
}
