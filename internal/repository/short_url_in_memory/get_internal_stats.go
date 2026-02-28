package shorturlinmemory

import (
	"context"
)

func (repo *InMemoryShortURLRepository) GetInternalStats(ctx context.Context) (urls int, users int, err error) {
	countURLs, countUsers := cont.countStats()

	return countURLs, countUsers, nil
}
