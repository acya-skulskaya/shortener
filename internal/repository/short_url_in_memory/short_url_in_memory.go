package shorturlinmemory

import (
	"sync"

	errorsInternal "github.com/acya-skulskaya/shortener/internal/errors"
)

type InMemoryShortURLRepository struct {
}

const containerMapSize = 250000

// generate:reset
type Container struct {
	mu        sync.RWMutex
	shortUrls map[string]shortURL
}

// generate:reset
type shortURL struct {
	shortURL    string
	originalURL string
	userID      string
	isDeleted   int8
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
	c.shortUrls[id] = shortURL{originalURL: value, userID: userID, shortURL: id, isDeleted: 0}

	return id, nil
}

func (c *Container) getItem(id string) (item shortURL, err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.shortUrls[id]
	if !ok {
		return shortURL{}, errorsInternal.ErrIDNotFound
	}

	return item, nil
}

func (c *Container) deleteItem(id string) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.shortUrls[id]
	if !ok {
		return errorsInternal.ErrIDNotFound
	}
	item.isDeleted = 1
	c.shortUrls[id] = item

	return nil
}

func (c *Container) getByUserID(userID string) (list []shortURL) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, item := range c.shortUrls {
		if item.isDeleted == 1 {
			continue
		}

		if item.userID == userID {
			list = append(list, item)
		}
	}

	return list
}

var cont = Container{shortUrls: make(map[string]shortURL, containerMapSize)}
