package vermis

import "errors"

var (
	ErrNotSetHostAddress = errors.New("don't set host address")
	ErrWal               = "wal error:: %w"
)

type Vermis interface {
	Add(el Element)
	Len() uint64
	GetFromIdx(uidx uint64) []Element
	SetHost() error
	SetReplica() error
	Stop()
}

type Element interface{}

type UnmarshalFunc func([]byte) (any, error)
