package syncmap

import "sync"

type Map[Tk any, Tv any] struct {
	m *sync.Map
}

func New[Tk any, Tv any]() *Map[Tk, Tv] {
	return &Map[Tk, Tv]{
		m: &sync.Map{},
	}
}

func (r *Map[Tk, Tv]) Store(name Tk, value Tv) {
	r.m.Store(name, value)
}

func (r *Map[Tk, Tv]) LoadOrStore(name Tk, value Tv) (v Tv, ok bool) {
	actual, loaded := r.m.LoadOrStore(name, value)
	return actual.(Tv), loaded
}

func (r *Map[Tk, Tv]) Load(name Tk) (v Tv) {
	load, ok := r.m.Load(name)
	if ok {
		v = load.(Tv)
	}
	return
}

func (r *Map[Tk, Tv]) Range(fn func(Tk, Tv) bool) {
	r.m.Range(func(key, value any) bool {
		return fn(key.(Tk), value.(Tv))
	})
}

func (r *Map[Tk, Tv]) LoadAndDelete(name Tk) (v Tv) {
	load, ok := r.m.LoadAndDelete(name)
	if ok {
		v = load.(Tv)
	}
	return
}

func (r *Map[Tk, Tv]) Delete(name Tk) {
	r.m.Delete(name)
}

func (r *Map[Tk, Tv]) Exist(name Tk) bool {
	_, ok := r.m.Load(name)
	return ok
}
