package vermis

import (
	"log"
	"os"
	"sync"

	"github.com/tidwall/wal"
)

type SimpleVermis struct {
	sync.RWMutex
	data       []Element
	file       *os.File
	wal        *wal.Log
	writerChan chan Element
	doneChan   chan struct{}
}

func NewSimpleVermis(filePath string, f UnmarshalFunc) (*SimpleVermis, error) {
	var err error

	s := SimpleVermis{
		writerChan: make(chan Element),
		doneChan:   make(chan struct{}),
	}

	s.wal, err = wal.Open(filePath, nil)
	if err != nil {
		return nil, err
	}
	if err = s.readWal(f); err != nil {
		return nil, err
	}

	go s.writerBG()

	return &s, nil
}

func (s *SimpleVermis) Add(el Element) {
	s.add(el)
	s.writerChan <- el
}

func (s *SimpleVermis) add(el Element) {
	s.Lock()
	defer s.Unlock()
	s.data = append(s.data, el)
}

func (s *SimpleVermis) Len() uint64 {
	s.RLock()
	defer s.RUnlock()

	return uint64(len(s.data))
}

func (s *SimpleVermis) Get(idx int64) []any {
	// TODO implement me
	panic("implement me")
}

func (s *SimpleVermis) SetMaster() error {
	// TODO implement me
	panic("implement me")
}

func (s *SimpleVermis) SetSlave(address string) error {
	// TODO implement me
	panic("implement me")
}

func (s *SimpleVermis) Stop() {
	if err := s.file.Close(); err != nil {
		log.Println(err.Error())
	}
	close(s.doneChan)
}

func (s *SimpleVermis) Show(prefix string) {
	s.Lock()
	defer s.Unlock()
	for i, r := range s.data {
		log.Printf(prefix+"::id=%d::%s\n", i, r)
	}
}
