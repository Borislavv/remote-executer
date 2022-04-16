package main

import "https://github.com/Borislavv/remote-executer/cmd/remoter"

func main() {
	if err := remoter.Run(); err != nil {
		log.Fatalln()
	}
}
