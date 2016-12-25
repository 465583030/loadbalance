package loadbalance

import "testing"

func BenchmarkSafeNginxScheduler_Next(b *testing.B) {
	s := NewSafeNginxScheduler(BuildNodes([]int{1, 2, 3}))
	for i := 0; i < b.N; i++ {
		s.Next()
	}
}

func BenchmarkSafeNginxScheduler_Next_Paral(b *testing.B) {
	s := NewSafeNginxScheduler(BuildNodes([]int{1, 2, 3}))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Next()
		}
	})
}
