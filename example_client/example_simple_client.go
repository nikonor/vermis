package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/nikonor/vermis/simplevermis"
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

	client, err := simplevermis.NewSimpleVermis(
		"/tmp/client.wal",
		"localhost:9990",
		9991,
		f,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	client.Show("begin")
	defer client.Show("end")

	if err = client.SetReplica(); err != nil {
		log.Fatalln(err)
	}

	time.Sleep(time.Second / 1000)
	fmt.Print("Enter numbers + Enter ( 0 - for exit ): ")
	for {
		var d int
		fmt.Scanf("%d", &d)
		if d == 0 {
			client.Stop()
			return
		} else {
			fmt.Println("")
		}
	}

	// if err := client.SetHost(); err != nil {
	// 	log.Fatalln(err.Error())
	// }
	//
	// time.Sleep(3 * time.Second)
	// client.Stop()
	// time.Sleep(time.Second)
}
