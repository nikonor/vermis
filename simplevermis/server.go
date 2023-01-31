package simplevermis

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/nikonor/vermis"
)

func (s *SimpleVermis) startServer() error {
	log.Println("call startServer")
	var (
		err  error
		port = ":" + strconv.Itoa(s.port)
	)

	if !s.iAmHost {
		return vermis.ErrNotHost
	}
	if s.server != nil {
		return vermis.ErrServerAlreadyStarted
	}

	if s.server, err = net.Listen("tcp", port); err != nil {
		return fmt.Errorf(vermis.ErrNetwork, "start listener", err)
	}

	for {
		var cli net.Conn
		cli, err = s.server.Accept()
		if err != nil {
			log.Println(err.Error())
		}

		go s.startHandler(cli)
	}

	log.Println("server started...")
	return nil
}

func (s *SimpleVermis) startHandler(cli net.Conn) {
	log.Println("start handler")
	defer log.Println("finish handler")

	for {
		netData, err := bufio.NewReader(cli).ReadBytes(vermis.TERMINATOR)
		if err != nil {
			log.Println(fmt.Errorf(vermis.ErrNetwork, "read", err))
			break
		}

		temp := bytes.TrimSpace(netData)
		fmt.Println("Received:", temp)
		cli.Write(vermis.CMDAck)
	}
}

func (s *SimpleVermis) stopServer() error {
	log.Println("server started...")

	var (
		err error
	)

	if err = s.server.Close(); err != nil {
		return fmt.Errorf(vermis.ErrNetwork, "close listener", err)
	}

	s.server = nil

	return nil
}
