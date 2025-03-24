package storage

import (
	"songsapi/query"
)

type Storage[T any] interface {
	Get(id int) (*T, error)
	Create(model *T) error
	Delete(model *T) error
	Update(model *T) error
	Find(q query.Query) ([]*T, error)
}
