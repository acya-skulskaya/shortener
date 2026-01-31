package pool

import (
	"sync"
)

type Resettable interface {
	Reset()
}

type ResettablePool[T Resettable] struct {
	pool sync.Pool
}

func New[T Resettable](newFunc func() T) *ResettablePool[T] {
	return &ResettablePool[T]{
		pool: sync.Pool{
			New: func() any {
				return newFunc()
			},
		},
	}
}

func (p *ResettablePool[T]) Get() T {
	return p.pool.Get().(T)
}

func (p *ResettablePool[T]) Put(r T) {
	r.Reset()
	p.pool.Put(r)
}
