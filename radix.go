package idempotent

import (
	iradix "github.com/hashicorp/go-immutable-radix"
)

type RadixTree struct {
	r *iradix.Tree
}

func (rt *RadixTree) Init() {
	rt.r = iradix.New()
	log.Info("radix tree inited.")

}

func (rt *RadixTree) Duplicated(key interface{}) (bool, error) {
	return false, nil
}

func (rt *RadixTree) Save(key interface{}) error {

	return nil
}
