package main

import (
	"log"

	"github.com/Borislavv/remote-executer/cmd/remoter"
)

func main() {
	// init. remoter-executer service
	if err := remoter.Run(); err != nil {
		log.Fatalln()
	}
}
