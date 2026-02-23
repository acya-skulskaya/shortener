package http

import (
	"github.com/acya-skulskaya/shortener/internal/observer/audit/publisher"
	interfaces "github.com/acya-skulskaya/shortener/internal/repository/interface"
)

// ShortUrlsService provides access to URL Shortener storage interface and audit publisher
type ShortUrlsService struct {
	Repo           interfaces.ShortURLRepository
	auditPublisher publisher.Publisher
}

// NewShortUrlsService creates a new instance of ShortUrlsService
func NewShortUrlsService(su interfaces.ShortURLRepository, ap publisher.Publisher) *ShortUrlsService {
	return &ShortUrlsService{
		Repo:           su,
		auditPublisher: ap,
	}
}
