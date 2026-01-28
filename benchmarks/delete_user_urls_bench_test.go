package main

import (
	"context"
	"fmt"
	shorturlinmemory "github.com/acya-skulskaya/shortener/internal/repository/short_url_in_memory"
	"testing"
)

func BenchmarkDeleteUserURLs(b *testing.B) {
	repo := shorturlinmemory.InMemoryShortURLRepository{}
	_, users, usersToURLs := seedShortUrls(&repo, 10000)

	b.ResetTimer()
	j := 0
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		if j >= len(users) {
			fmt.Println(i)
			j = 0
			_, users, usersToURLs = seedShortUrls(&repo, 10000)
		}
		b.StartTimer()
		repo.DeleteUserUrls(context.Background(), usersToURLs[users[j]], users[j])
		j++
	}
}
