package simplevermis

import (
	"log"
	"sync"

	"github.com/nikonor/vermis"
	"github.com/tidwall/wal"
)

type SimpleVermis struct {
	sync.RWMutex
	data        []vermis.Element
	wal         *wal.Log
	writerChan  chan vermis.Element
	doneChan    chan struct{}
	iAmHost     bool
	hostAddress string
}

func NewSimpleVermis(filePath string, hostAddress string, f vermis.UnmarshalFunc) (*SimpleVermis, error) {
	var err error

	s := SimpleVermis{
		writerChan: make(chan vermis.Element),
		doneChan:   make(chan struct{}), hostAddress: hostAddress,
	}

	if len(hostAddress) == 0 {
		return nil, vermis.ErrNotSetMasterAddress
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

func (s *SimpleVermis) Add(el vermis.Element) {
	s.add(el)
	s.writerChan <- el
}

func (s *SimpleVermis) add(el vermis.Element) {
	s.Lock()
	defer s.Unlock()
	s.data = append(s.data, el)
}

func (s *SimpleVermis) Len() uint64 {
	s.RLock()
	defer s.RUnlock()

	return uint64(len(s.data))
}

func (s *SimpleVermis) GetFromIdx(idx uint64) []vermis.Element {
	s.Lock()
	defer s.Unlock()
	return s.data[int(idx):]
}

func (s *SimpleVermis) SetHost() error {
	log.Println("set instance as host")
	s.iAmHost = true
	return nil

}

func (s *SimpleVermis) SetReplica() error {
	log.Println("set instance as replica")
	s.iAmHost = false
	return nil
}

func (s *SimpleVermis) Stop() {
	if err := s.wal.Sync(); err != nil {
		log.Println(err.Error())
	}

	if err := s.wal.Close(); err != nil {
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
