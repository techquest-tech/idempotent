package idempotent

import "github.com/sirupsen/logrus"

var log = logrus.WithField("component", "idempotent")

//IdempotentKey key interface
type IdempotentKey interface {
	IdempotentKey() (interface{}, error)
}

//IdempotentService Idempotent Service
type IdempotentService interface {
	//Duplicated Idempotent checking, true if it's duplicated request.
	Duplicated(key interface{}) (bool, error)
	Save(key interface{}) error
	AllKeys() ([]interface{}, error)
}
