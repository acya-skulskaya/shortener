package interfaces

import (
	"context"

	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
)

// ShortURLRepository defines an interface to operate on shortened URLs
// It provides methods to create, retrieve, and manage shortened URLs.
type ShortURLRepository interface {

	// GetInternalStats retrieves number of users and shortened URLs
	GetInternalStats(ctx context.Context) (urls int, users int, err error)

	// Get retrieves a shortened URL by ID or returns an error if ID does not exist
	Get(ctx context.Context, id string) (originalURL string, err error)

	// GetUserUrls retrieves all shortened URLs that were created by a specified user
	GetUserUrls(ctx context.Context, userID string) (list []jsonModel.BatchURLList, err error)

	// DeleteUserUrls deletes all shortened URLs that were created by a specified user
	DeleteUserUrls(ctx context.Context, list []string, userID string)

	// Store creates a shortened URL and returns its ID
	Store(ctx context.Context, originalURL string, userID string) (id string, err error)

	// StoreBatch creates a list of shortened URL and returns their ID
	StoreBatch(ctx context.Context, listOriginal []jsonModel.BatchURLList, userID string) (listShorten []jsonModel.BatchURLList, err error)
}
