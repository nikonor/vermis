package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"time"
)

type Message struct {
	Msg string
}

func main() {
	gob.Register(new(Message))

	serverListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println(err)
	}
	conn, err := net.Dial("tcp", serverListener.Addr().String())
	if err != nil {
		fmt.Println(err)
	}
	serverConn, err := serverListener.Accept()
	if err != nil {
		fmt.Println(err)
	}
	done := false
	dec := gob.NewDecoder(conn) // Will read from network.
	enc := gob.NewEncoder(serverConn)
	go func() {
		for !done {
			recieveMessage(dec)
		}
	}()

	for i := 1; i < 1000; i++ {
		sent := Message{strconv.Itoa(i)}
		sendMessage(sent, enc)
	}
	time.Sleep(time.Second)
	done = true
}

func sendMessage(msg Message, enc *gob.Encoder) {
	err := enc.Encode(msg)
	if err != nil {
		fmt.Println(err)
	}
}

func recieveMessage(dec *gob.Decoder) {
	msg := new(Message)
	err := dec.Decode(msg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Client recieved:", msg.Msg)
}
