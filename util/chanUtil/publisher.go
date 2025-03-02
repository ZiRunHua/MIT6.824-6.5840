package chanUtil

import (
	"context"
)

type (
	Publisher[T any] interface {
		Publish(t T)
	}
	publisher[T any] struct {
		onPublish func(t T)
		queue     chan T
		ctx       context.Context
	}
)

func NewPublisher[T any](size int, onPublish func(t T), ctx context.Context) Publisher[T] {
	p := &publisher[T]{
		onPublish: onPublish,
		queue:     make(chan T, size),
		ctx:       ctx,
	}
	go func() {
		for {
			select {
			case t := <-p.queue:
				p.onPublish(t)
			case <-p.ctx.Done():
				return
			}
		}
	}()
	return p
}

func (p *publisher[T]) Publish(t T) {
	select {
	case p.queue <- t:
	case <-p.ctx.Done():
		return
	}
}
