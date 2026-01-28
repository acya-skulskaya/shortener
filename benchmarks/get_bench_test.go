package main

//func BenchmarkGet(b *testing.B) {
//	repo := shorturlinmemory.InMemoryShortURLRepository{}
//	urlIDs, _, _ := seedShortUrls(&repo, 10000)
//
//	b.ResetTimer()
//	j := 0
//	for i := 0; i < b.N; i++ {
//		b.StopTimer()
//		if j >= len(urlIDs) {
//			j = 0
//		}
//		b.StartTimer()
//		repo.Get(context.Background(), urlIDs[j])
//		j++
//	}
//}
