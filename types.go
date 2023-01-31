package vermis

import "errors"

const (
	TERMINATOR = '@'
)

var (
	CMDAck = []byte("ack@")

	ErrNotSetHostAddress    = errors.New("don't set host address")
	ErrNotSetPort           = errors.New("don't set port")
	ErrNotHost              = errors.New("instance is not host")
	ErrServerAlreadyStarted = errors.New("server already started")
	ErrWal                  = "wal error::%s::%w"
	ErrNetwork              = "network error::%s::%w"
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
