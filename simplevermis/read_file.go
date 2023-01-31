package simplevermis

import (
	"bufio"
	"log"
	"os"

	"github.com/nikonor/vermis"
)

func (s *SimpleVermis) readByLines(file *os.File, f vermis.UnmarshalFunc) error {

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		tmp := scanner.Text()
		if e, err := f([]byte(tmp)); err != nil {
			log.Println("\terr=", err.Error())
		} else {
			s.add(e)
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (s *SimpleVermis) readWal(f vermis.UnmarshalFunc) error {
	l, err := s.wal.LastIndex()
	if err != nil {
		return err
	}
	for i := uint64(1); i <= l; i++ {
		body, err := s.wal.Read(i)
		if err != nil {
			return err
		}
		if e, err := f(body); err != nil {
			return err
		} else {
			s.add(e)
		}

	}

	return nil
}
