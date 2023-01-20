package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/nikonor/vermis"
)

type E struct {
	One string
	Two int
}

func (e E) String() string {
	return "one=" + e.One + ",two=" + strconv.Itoa(e.Two)
}

func main() {
	log.Println("Start")
	defer log.Println("Finish")

	f := func(a []byte) (any, error) {
		var e E
		if err := json.Unmarshal(a, &e); err != nil {
			return E{}, err
		}
		return e, nil
	}

	srv, err := vermis.NewSimpleVermis("/tmp/server.txt", f)
	if err != nil {
		log.Fatalln(err.Error())
	}

	srv.Show("begin")
	defer srv.Show("end")

	for {
		var d int
		fmt.Scanf("%d", &d)
		if d == 0 {
			srv.Stop()
			return
		}

		srv.Add(E{
			One: time.Now().Format("2006-01-02 15:04:05"),
			Two: d,
		})

	}

	// if err := srv.SetMaster(); err != nil {
	// 	log.Fatalln(err.Error())
	// }
	//
	// time.Sleep(3 * time.Second)
	// srv.Stop()
	// time.Sleep(time.Second)
}
