package idempotent

import (
	"fmt"
	"reflect"
	"sync"
)

// type IdempotentType uint

// const (
// 	InMem IdempotentType = iota
// )

type Idempotent struct {
	mu      sync.RWMutex
	Key     *DefaultIdempotentKey
	Service IdempotentService
}

func NewIdempotentWithTemplate(template string, service IdempotentService) (*Idempotent, error) {
	// service := service
	key, err := TemplateAsKey(template)
	if err != nil {
		return nil, err
	}
	return &Idempotent{
		Key:     key,
		Service: service,
	}, nil
}

func NewIdempotent(service IdempotentService) (*Idempotent, error) {
	return &Idempotent{
		Service: service,
	}, nil
}

// func GetIdempotent(t IdempotentType) *Idempotent {
// 	service := NewInMemoryMap()
// 	return &Idempotent{
// 		Service: service,
// 	}
// }

func (factory *Idempotent) Duplicated(obj interface{}) (bool, error) {
	if obj == nil {
		return true, nil
	}

	factory.mu.Lock()
	defer factory.mu.Unlock()

	var id interface{}
	var err error

	switch reflect.TypeOf(obj).Kind() {

	case reflect.Array, reflect.Func, reflect.Chan, reflect.Slice:
		err = fmt.Errorf("%s is not supported", reflect.TypeOf(obj).Kind())

	case reflect.Struct, reflect.Map:

		idObj, ok := obj.(IdempotentKey)
		if !ok {
			if factory.Key == nil {
				return true, fmt.Errorf("failed to get key from object, %T not imple IdempotentKey", obj)
			}

			factory.Key.Target = obj
			idObj = factory.Key
			log.Debugf("user default IdempotentKey for %T", obj)
		}

		id, err = idObj.IdempotentKey()

	default:
		id = obj
	}

	if id == nil {
		log.Warnf("key is nil for %+v", obj)
		return true, nil
	}

	if err != nil {
		log.Error("failed to get key from object, err ", err)
		return true, err
	}
	duplicated, err := factory.Service.Duplicated(id)
	if err != nil {
		log.Error("check duplicated failed. error: ", err)
		return true, err
	}
	if !duplicated {
		log.Debug("new key found ", id)
		err = factory.Service.Save(id)
		if err != nil {
			log.Errorf("Idempotent check done, but failed when save key %s. err %v ", id, err)
		}
	}
	return duplicated, err
}
