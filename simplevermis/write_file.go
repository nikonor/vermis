package simplevermis

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nikonor/vermis"
)

func (s *SimpleVermis) writerBG() {
	log.Println("start writerBG")
	for {
		select {
		case el := <-s.writerChan:
			// TODO: подумать над s.wal.WriteBatch()
			log.Println("writerBG got new element::", el)

			b, _ := json.Marshal(el)
			b = append(b, []byte("\n")...)

			if err := s.wal.Write(s.Len(), b); err != nil {
				log.Println(fmt.Errorf(vermis.ErrWal, err))
			}

		case <-s.doneChan:

			log.Println("finish writerBG")
			return

		}
	}
}
