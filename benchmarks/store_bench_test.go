package main

import (
	"context"
	shorturlinmemory "github.com/acya-skulskaya/shortener/internal/repository/short_url_in_memory"
	"github.com/google/uuid"
	"testing"
)

func BenchmarkStore(b *testing.B) {
	repo := shorturlinmemory.InMemoryShortURLRepository{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.Store(context.Background(), "http://test.test/"+uuid.New().String(), uuid.New().String())
	}
}
