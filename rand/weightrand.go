package rand

import (
	"sync"
)

type SeedType interface {
	any
}

type Item[T SeedType] struct {
	item   T
	weight int
}

type WeightRand[T SeedType] struct {
	seed   []Item[T]
	weight int
	lock   *sync.Mutex
}

func NewWeightRand[T SeedType](items ...T) *WeightRand[T] {
	w := &WeightRand[T]{
		lock: &sync.Mutex{},
	}

	for _, it := range items {
		w.Add(it)
	}

	return w
}

func (w *WeightRand[T]) Add(it T) {
	w.AddWeight(it, 1)
}

func (w *WeightRand[T]) AddWeight(it T, weight int) {
	defer w.calc()

	w.seed = append(w.seed, Item[T]{
		item:   it,
		weight: weight,
	})
}

func (w *WeightRand[T]) calc() {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.weight = 0
	for _, it := range w.seed {
		w.weight += it.weight
	}
}

func (w *WeightRand[T]) Get() (t T) {
	stop := Intn(w.weight)
	sum := 0
	for _, it := range w.seed {
		if sum += it.weight; sum > stop {
			return it.item
		}
	}

	return
}
