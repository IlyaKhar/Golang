package storage

import (
	"errors"
)

var ErrNotFound = errors.New("key not found")

type Storage interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	Has(key string) bool
	Keys() []string
	Len() int
	Clear() error
}
