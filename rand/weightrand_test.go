package rand

import (
	"fmt"
	"sync"
	"testing"
)

func TestWeightRand(t *testing.T) {
	r := NewWeightRand[int]()
	r.AddWeight(14, 0)

	l := 10
	highWeightItem := 0
	m := map[int]int{}
	for i := 0; i < l; i++ {
		if i == highWeightItem {
			r.AddWeight(i, 2)
		} else {
			r.Add(i)
		}
	}

	count := 1000000
	for i := 0; i < count; i++ {
		m[r.Get()]++
	}

	if _, ok := m[14]; ok || m[highWeightItem] <= count/l {
		t.Errorf("test failed")
	}
}

// go test -bench . -benchmem -run none -count 10 -memprofile mem.out
// go tool pprof -svg mem.out > mem.svg
func BenchmarkTestWeightGet(b *testing.B) {
	r := NewWeightRand[string]()

	wg := sync.WaitGroup{}
	for i := 0; i < 80000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r.Add(fmt.Sprintf("rand-it-%d", i))
		}(i)
	}

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r.AddWeight(fmt.Sprintf("rand-it-%d", i), 100)
		}(i)
	}

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r.AddWeight(fmt.Sprintf("rand-it-%d", i), 30)
		}(i)
	}

	wg.Wait()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Get()
	}
}
