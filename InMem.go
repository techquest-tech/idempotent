package idempotent

type InMemoryMap struct {
	// mu     sync.RWMutex
	NoLock bool
	cache  map[interface{}]bool
}

func NewInMemoryMap() *InMemoryMap {
	return &InMemoryMap{
		cache: map[interface{}]bool{},
	}
}

// func NewInMemoryMapWithoutLocker() *InMemoryMap {
// 	return &InMemoryMap{
// 		NoLock: true,
// 		cache:  map[interface{}]bool{},
// 	}
// }

func (mem *InMemoryMap) Duplicated(key interface{}) (bool, error) {

	_, ok := mem.cache[key]
	return ok, nil
}

func (mem *InMemoryMap) Save(key interface{}) error {
	if key == nil {
		log.Warn("key is nil, ignored")
		return nil
	}

	mem.cache[key] = true
	return nil
}

func (mem *InMemoryMap) AllKeys() ([]interface{}, error) {
	all := make([]interface{}, 0)

	for k := range mem.cache {
		all = append(all, k)
	}

	return all, nil
}
