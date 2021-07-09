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
	return ok, nil
}

func (mem *InMemoryMap) Save(key interface{}) error {
	if key == nil {
		log.Warn("key is nil, ignored")
		return nil
	}
	mem.mu.Lock()
	defer mem.mu.Unlock()
	mem.cache[key] = true
	return nil
}

func (mem *InMemoryMap) AllKeys() ([]interface{}, error) {
	all := make([]interface{}, 0)

	for k, _ := range mem.cache {
		all = append(all, k)
	}

	return all, nil
}
