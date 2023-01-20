package vermis

import (
	"bufio"
	"log"
	"os"
)

func (s *SimpleVermis) readByLines(file *os.File, f UnmarshalFunc) error {

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
