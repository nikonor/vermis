package main

import (
	"log"
)

func main() {
	log.Println("Start")
	defer log.Println("Finish")
}
