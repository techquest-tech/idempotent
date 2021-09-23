package idempotent

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	// "github.com/prometheus/common/log"
	"go.uber.org/zap"
)

// type IdempotentType uint

// const (
// 	InMem IdempotentType = iota
// )

type Idempotent struct {
	mu      sync.RWMutex
	Key     *DefaultIdempotentKey
	Service IdempotentService
	Logger  *zap.Logger
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

func (factory *Idempotent) GetObjectKey(obj interface{}) (interface{}, error) {
	if obj == nil {
		return nil, errors.New("object is nil")
	}

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
			// log.Debugf("user default IdempotentKey for %T", obj)
			factory.Logger.Debug("user default idempotent key", zap.Any("target", obj))
		}

		id, err = idObj.IdempotentKey()

	default:
		id = obj
	}

	if id == nil {
		// log.Warnf("key is nil for %+v", obj)
		factory.Logger.Warn("key is nil", zap.Any("target", obj))
		return true, nil
	}

	if err != nil {
		factory.Logger.Error("failed to get key from object, err ", zap.Any("error", err))
		return nil, err
	}
	return id, err
}

func (factory *Idempotent) Duplicated(obj interface{}) (bool, error) {

	id, err := factory.GetObjectKey(obj)
	if err != nil {
		return true, err
	}
	factory.mu.Lock()
	defer factory.mu.Unlock()

	duplicated, err := factory.Service.Duplicated(id)
	if err != nil {
		factory.Logger.Error("check duplicated failed. error: ", zap.Any("error", err))
		return true, err
	}
	if !duplicated {
		factory.Logger.Debug("new key found ", zap.Any("id", id))
		err = factory.Service.Save(id)
		if err != nil {
			factory.Logger.Error("Idempotent check done, but failed when save key", zap.Any("id", id), zap.Any("error", err))
		}
	}
	return duplicated, err
}
