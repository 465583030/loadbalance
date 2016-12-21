package loadbalance

import (
	"reflect"
	"testing"
)

func TestNginxScheduler_Next(t *testing.T) {
	cases := []struct {
		weights []int
		list    []int
	}{
		{
			[]int{1, 1, 1},
			[]int{0, 1, 2, 0, 1, 2},
		},
		{
			[]int{1, 2, 3},
			[]int{2, 1, 0, 2, 1, 2, 2, 1, 0, 2, 1, 2},
		},
		{
			[]int{2, 2, 4},
			[]int{2, 0, 1, 2, 2, 0, 1, 2},
		},
		{
			[]int{0, 2, 3},
			[]int{2, 1, 0, 2, 1, 2, 2, 1, 0, 2},
		},
		{
			[]int{-3, 2, 4},
			[]int{2, 0, 1, 2, 0, 2, 1, 0, 2, 2, 0, 1},
		},
		{
			[]int{3},
			[]int{0, 0},
		},
	}
	for _, c := range cases {
		scheduler := NewNginxScheduler(BuildNodes(c.weights))
		list := make([]int, len(c.list))
		for i, n := 0, len(c.list); i < n; i++ {
			list[i] = scheduler.Next().Data.(int)
		}
		if !reflect.DeepEqual(list, c.list) {
			t.Errorf("weights: %v, list: %v, want %v", c.weights, list, c.list)
		}
	}

	s := NewNginxScheduler(BuildNodes([]int{}))
	node := s.Next()
	if node != nil {
		t.Errorf("weights: []int{}, return %v, want nil", node)
	}
}

func BenchmarkNginxScheduler_Next(b *testing.B) {
	s := NewNginxScheduler(BuildNodes([]int{1, 2, 3}))
	for i := 0; i < b.N; i++ {
		s.Next()
	}
}

