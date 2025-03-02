package dataTool

import "sync"

type RWMutexMap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

func NewRWMutexMap[K comparable, V any]() *RWMutexMap[K, V] {
	return &RWMutexMap[K, V]{m: make(map[K]V)}
}

func (sm *RWMutexMap[K, V]) Load(key K) (value V, ok bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	value, ok = sm.m[key]
	return
}

func (sm *RWMutexMap[K, V]) LoadOrZero(key K) (value V) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	value, _ = sm.m[key]
	return
}

func (sm *RWMutexMap[K, V]) Store(key K, value V) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m[key] = value
}

func (sm *RWMutexMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if actual, loaded = sm.m[key]; loaded {
		return actual, true
	}

	sm.m[key] = value
	return value, false
}

func (sm *RWMutexMap[K, V]) Delete(key K) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.m, key)
}

func (sm *RWMutexMap[K, V]) Range(f func(key K, value V) (shouldContinue bool)) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	for k, v := range sm.m {
		if !f(k, v) {
			break
		}
	}
}
