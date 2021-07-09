package idempotent

import (
	"sync"
)

type InMemoryMap struct {
	mu    sync.RWMutex
	cache map[interface{}]bool
}

func NewInMemoryMap() *InMemoryMap {
	return &InMemoryMap{
		cache: map[interface{}]bool{},
	}
}

func (mem *InMemoryMap) Duplicated(key interface{}) (bool, error) {
	mem.mu.RLock()
	defer mem.mu.RUnlock()

	_, ok := mem.cache[key]
	if !ok {
		mem.cache[key] = true
	}
	return ok, nil
}

func (mem *InMemoryMap) Save(key interface{}) error {
	mem.mu.Lock()
	defer mem.mu.Unlock()
	mem.cache[key] = true
	return nil
}
