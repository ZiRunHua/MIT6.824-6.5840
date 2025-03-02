package chanUtil

import "sync"

type (
	Broadcaster interface {
		Broadcast()
		Close()
	}
	broadcast struct {
		channels []chan struct{}

		done chan struct{}
		mu   sync.Mutex
	}
)

func NewBroadcaster(callbacks []func()) Broadcaster {
	p := &broadcast{
		channels: make([]chan struct{}, len(callbacks), len(callbacks)),
		done:     make(chan struct{}),
	}

	for i, callback := range callbacks {
		p.channels[i] = make(chan struct{}, 1)
		go listener(p.channels[i], p.done, callback)
	}
	return p
}

func (p *broadcast) Close() {
	defer func() { _ = recover() }()
	close(p.done)
}

func (p *broadcast) Broadcast() {
	for _, c := range p.channels {
		select {
		case <-p.done:
			return
		case c <- struct{}{}:
		default:
		}
	}
}
