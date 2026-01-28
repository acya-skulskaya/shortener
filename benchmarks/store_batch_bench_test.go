package main

import (
	"context"
	jsonModel "github.com/acya-skulskaya/shortener/internal/model/json"
	shorturlinmemory "github.com/acya-skulskaya/shortener/internal/repository/short_url_in_memory"
	"github.com/google/uuid"
	"testing"
)

func BenchmarkStoreBatch(b *testing.B) {
	repo := shorturlinmemory.InMemoryShortURLRepository{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.StoreBatch(context.Background(), []jsonModel.BatchURLList{
			{
				CorrelationID: uuid.New().String(),
				OriginalURL:   "http://test.test/" + uuid.New().String(),
			}, {
				CorrelationID: uuid.New().String(),
				OriginalURL:   "http://test.test/" + uuid.New().String(),
			}, {
				CorrelationID: uuid.New().String(),
				OriginalURL:   "http://test.test/" + uuid.New().String(),
			}, {
				CorrelationID: uuid.New().String(),
				OriginalURL:   "http://test.test/" + uuid.New().String(),
			},
		}, uuid.New().String())
	}
}
