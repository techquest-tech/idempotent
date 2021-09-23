package idempotent

import (
	iradix "github.com/armon/go-radix"
)

type RadixTree struct {
	r *iradix.Tree
}

func NewRadixTree() *RadixTree {
	t := &RadixTree{}
	t.Init()
	return t
}

func (rt *RadixTree) Init() {
	rt.r = iradix.New()
	// rt.Logger.Info("radix tree inited.")

}

func (rt *RadixTree) Duplicated(key interface{}) (bool, error) {
	btKey, ok := key.(string)
	if !ok {
		return true, errTypeNotSupport
	}
	// if rt.r.Root() == nil {
	// 	return false, nil
	// }
	_, ok = rt.r.Get(btKey)
	return ok, nil
}

func (rt *RadixTree) Save(key interface{}) error {
	btKey, ok := key.(string)
	if !ok {
		return errTypeNotSupport
	}
	rt.r.Insert(btKey, true)
	// ntree, _, _ := rt.r.Insert(btKey, true)
	// rt.r = ntree
	// if !ok {
	// 	log.Warnf("save key %s return false", btKey)
	// }
	return nil
}

func (rt *RadixTree) AllKeys() ([]interface{}, error) {
	results := make([]interface{}, 0)
	walk := func(k string, v interface{}) bool {
		results = append(results, k)
		return false
	}
	rt.r.Walk(walk)
	return results, nil
}
