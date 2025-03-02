package chanUtil

import (
	"context"
	"sync"
)

type (
	// Throttler 节流器 用于异步调用，Publish 时不会堵塞调用方
	Throttler interface {
		Publish()
		Close()
	}
	throttle struct {
		callback func()

		done, notifyCh chan struct{}
		mu             sync.Mutex
	}
)

func NewThrottle(callback func()) Throttler {
	p := &throttle{
		callback: callback,
		done:     make(chan struct{}),
		notifyCh: make(chan struct{}, 1),
	}
	go listener(p.notifyCh, p.done, callback)
	return p
}

func (p *throttle) Close() {
	defer func() { _ = recover() }()
	if p == nil {
		return
	}
	close(p.done)
}

func (p *throttle) Publish() {
	if p == nil {
		return
	}
	select {
	case <-p.done:
		return
	case p.notifyCh <- struct{}{}:
	default:
	}
}

type (
	// BlockingThrottler 会堵塞调用方的节流器 适合用在lab3的持久化中
	BlockingThrottler interface {
		Run()
	}
	blockingThrottler struct {
		execute func()
		mu      sync.Mutex
		ch      chan struct{}
		running bool

		ctx context.Context
	}
)

func NewBlockingThrottler(f func(), ctx context.Context) BlockingThrottler {
	return &blockingThrottler{
		execute: f,
		ch:      make(chan struct{}, 1),
		ctx:     ctx,
	}
}

func (t *blockingThrottler) Run() {
	t.mu.Lock()
	if t.running {
		// 如果正在运行则等待当前运行结束
		t.mu.Unlock()
		if !t.waitRun() {
			return
		}
		t.mu.Lock()
	}
	if t.running {
		// 运行结束后 再次进入运行状态 说明其他协程触发了
		t.mu.Unlock()
		if !t.waitRun() {
			return
		}
	} else {
		// 运行结束后 没有再次进入运行状态 则由当前协程进行触发
		t.ch, t.running = make(chan struct{}, 1), true
		t.mu.Unlock()
		t.execute()
		t.mu.Lock()
		close(t.ch)
		t.running = false
		t.mu.Unlock()
	}
}
func (t *blockingThrottler) waitRun() bool {
	select {
	case <-t.ctx.Done():
		return false
	case <-t.ch:
		return true
	}
}
