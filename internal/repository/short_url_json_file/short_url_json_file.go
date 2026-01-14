package shorturljsonfile

import "sync"

type JSONFileShortURLRepository struct {
	FileStoragePath string
	mu              sync.RWMutex
}
