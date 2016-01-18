package loadbalance

import (
	"reflect"
	"testing"
)

func TestGcd(t *testing.T) {
	cases := []struct {
		a, b int
		gcd  int
	}{
		{0, 1, 1},
		{1, 0, 1},
		{1, 1, 1},
		{2, 2, 2},
		{4, 2, 2},
		{15, 2, 1},
		{13, 7, 1},
		{18, 9, 9},
		{15, 20, 5},
		{-10, 2, 1},
	}

	scheduler := NewWeightedRoundRobinScheduler(nil)
	for _, c := range cases {
		gcd := scheduler.gcd(c.a, c.b)
		if gcd != c.gcd {
			t.Errorf("gcd(%d, %d) == %d != %d", c.a, c.b, gcd, c.gcd)
		}
	}
}

func TestGcdSlice(t *testing.T) {
	cases := []struct {
		slice []int
		gcd   int
	}{
		{[]int{2, 4, 6, 8}, 2},
		{[]int{-3, 3, 6}, 1},
		{[]int{0, 4, 8}, 1},
		{[]int{1, 2, 4, 5}, 1},
		{[]int{10, 30, 50}, 10},
		{[]int{20, 40, 60}, 20},
	}
	scheduler := NewWeightedRoundRobinScheduler(nil)
	for _, c := range cases {
		gcd := scheduler.gcdSlice(c.slice)
		if gcd != c.gcd {
			t.Errorf("gcdSlice(%v) == %d != %d", c.slice, gcd, c.gcd)
		}
	}
}

func TestNext(t *testing.T) {
	cases := []struct {
		weights []int
		list    []int // 循环两次的索引列表
	}{
		{
			[]int{1, 1, 1},
			[]int{0, 1, 2, 0, 1, 2},
		},
		{
			[]int{1, 2, 3},
			[]int{2, 1, 2, 0, 1, 2, 2, 1, 2, 0, 1, 2},
		},
		{
			[]int{2, 2, 4},
			[]int{2, 0, 1, 2, 2, 0, 1, 2},
		},
		{
			[]int{0, 2, 3},
			[]int{2, 1, 2, 1, 2, 2, 1, 2, 1, 2},
		},
		{
			[]int{-3, 2, 4},
			[]int{2, 2, 1, 2, 1, 2, 2, 2, 1, 2, 1, 2},
		},
		{
			[]int{3},
			[]int{0, 0},
		},
	}
	for _, c := range cases {
		scheduler := NewWeightedRoundRobinScheduler(c.weights)
		list := make([]int, len(c.list))
		for i, n := 0, len(c.list); i < n; i++ {
			list[i] = scheduler.Next()
		}
		if !reflect.DeepEqual(list, c.list) {
			t.Errorf("weights: %v, list: %v, want %v", c.weights, list, c.list)
		}
	}
}
