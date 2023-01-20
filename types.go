package vermis

type Vermis interface {
	Add(el any)
	Len() int64
	Get(idx int64) []any
	SetMaster() error
	SetSlave(address string) error
	Stop()
}

type Element interface{}

type UnmarshalFunc func([]byte) (any, error)
