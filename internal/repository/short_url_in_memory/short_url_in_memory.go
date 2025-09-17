package short_url_in_memory

import (
	"github.com/acya-skulskaya/shortener/internal/helpers"
	"sync"
)

type InMemoryShortURLRepository struct {
}

type Container struct {
	mu        sync.Mutex
	shortUrls map[string]string
}

func (c *Container) add(id string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.shortUrls[id] = value
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

func (repo *InMemoryShortURLRepository) Store(originalURL string) (id string) {
	id = helpers.RandStringRunes(10)

	cont.add(id, originalURL)

	return id
}
