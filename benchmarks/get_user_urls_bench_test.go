package main

import (
	"context"
	shorturlinmemory "github.com/acya-skulskaya/shortener/internal/repository/short_url_in_memory"
	"testing"
)

func BenchmarkGetUserURLs(b *testing.B) {
	repo := shorturlinmemory.InMemoryShortURLRepository{}
	_, users, _ := seedShortUrls(&repo, 10000)

	b.ResetTimer()
	j := 0
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		if j >= len(users) {
			j = 0
		}
		b.StartTimer()
		repo.GetUserUrls(context.Background(), users[j])
		j++
	}
}
