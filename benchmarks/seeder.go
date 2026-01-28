package main

import (
	"context"
	"math/rand"

	interfaces "github.com/acya-skulskaya/shortener/internal/repository/interface"
	"github.com/google/uuid"
)

func randRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func seedShortUrls(repo interfaces.ShortURLRepository, count int) ([]string, []string, map[string][]string) {
	users := make([]string, count)
	urlIDs := make([]string, count)
	usersToURLs := make(map[string][]string, count)

	for i := 0; i < count; i++ {
		userID := uuid.New().String()

		users = append(users, userID)

		id, _ := repo.Store(context.Background(), "http://test.test/"+uuid.New().String(), userID)

		urlIDs = append(urlIDs, id)
		usersToURLs[userID] = []string{id}
	}

	return urlIDs, users, usersToURLs
}
