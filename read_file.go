package vermis

import (
	"bufio"
	"log"
	"os"
)

const bufferSize = 128

func (s *SimpleVermis) readByLines(file *os.File, f func([]byte) (any, error)) error {

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
