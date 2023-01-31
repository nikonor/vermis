package simplevermis

import (
	"fmt"
	"log"
	"net"
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
	port        int
	server      net.Listener
}

// TODO: заменить логи на вызов логировщика??? или вообще убрать

func NewSimpleVermis(filePath string, hostAddress string, port int, f vermis.UnmarshalFunc) (*SimpleVermis, error) {
	var err error

	s := SimpleVermis{
		writerChan:  make(chan vermis.Element),
		doneChan:    make(chan struct{}),
		hostAddress: hostAddress,
		port:        port,
	}

	if s.port <= 0 {
		return nil, vermis.ErrNotSetPort
	}
	if len(hostAddress) == 0 {
		return nil, vermis.ErrNotSetHostAddress
	}

	s.wal, err = wal.Open(filePath, nil)
	if err != nil {
		return nil, fmt.Errorf(vermis.ErrWal, "open", err)
	}
	if err = s.readWal(f); err != nil {
		return nil, fmt.Errorf(vermis.ErrWal, "read", err)
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
	return s.startServer()
}

func (s *SimpleVermis) SetReplica() error {
	log.Println("set instance as replica")

	if s.iAmHost {
		log.Println("instance was host. stop server")
		if err := s.stopServer(); err != nil {
			return err
		}
	}

	s.iAmHost = false

	return nil
}

func (s *SimpleVermis) Stop() {
	if s.iAmHost {
		if err := s.stopServer(); err != nil {
			log.Println(err.Error())
		}
	}

	if err := s.wal.Sync(); err != nil {
		log.Println(fmt.Errorf(vermis.ErrWal, "sync", err))
	}

	if err := s.wal.Close(); err != nil {
		log.Println(fmt.Errorf(vermis.ErrWal, "close", err))
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
