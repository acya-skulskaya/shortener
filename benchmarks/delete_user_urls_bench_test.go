package main

//func BenchmarkDeleteUserURLs(b *testing.B) {
//	repo := shorturlinmemory.InMemoryShortURLRepository{}
//	_, users, usersToURLs := seedShortUrls(&repo, 10000)
//
//	b.ResetTimer()
//	j := 0
//	for i := 0; i < b.N; i++ {
//		b.StopTimer()
//		if j >= len(users) {
//			j = 0
//			_, users, usersToURLs = seedShortUrls(&repo, 10000)
//		}
//		b.StartTimer()
//		repo.DeleteUserUrls(context.Background(), usersToURLs[users[j]], users[j])
//		j++
//	}
//}
