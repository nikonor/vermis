package vermis

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

type SimpleVermis struct {
	sync.RWMutex
	data       []Element
	file       *os.File
	writerChan chan Element
	doneChan   chan struct{}
}

func NewSimpleVermis(filePath string, f UnmarshalFunc) (*SimpleVermis, error) {
	var err error

	s := SimpleVermis{
		writerChan: make(chan Element),
		doneChan:   make(chan struct{}),
	}

	if s.file, err = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0777); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if s.file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0777); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if err = s.readByLines(s.file, f); err != nil {
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

func (s *SimpleVermis) Len() int64 {
	s.RLock()
	defer s.RUnlock()

	return int64(len(s.data))
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

func (s *SimpleVermis) writerBG() {
	log.Println("start writerBG")
	for {
		select {
		case el := <-s.writerChan:
			log.Println("writerBG got new element::", el)
			b, _ := json.Marshal(el)
			b = append(b, []byte("\n")...)
			if _, err := s.file.Write(b); err != nil {
				log.Println("error on write::", err)
			}
		case <-s.doneChan:
			log.Println("finish writerBG")
			return
		}
	}
}

func (s *SimpleVermis) Show(prefix string) {
	s.Lock()
	defer s.Unlock()
	for i, r := range s.data {
		log.Printf(prefix+"::id=%d::%s\n", i, r)
	}
}
