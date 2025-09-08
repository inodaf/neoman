package main

import (
	"fmt"
	"os"

	"github.com/inodaf/neoman/internal3/management"
	"github.com/inodaf/neoman/internal3/operations"
)

func main() {
	err := management.SocketClientPing()
	if err != nil {
		fmt.Println("neoman: Could not connect to daemon")
		return
	}

	if len(os.Args) == 1 || (len(os.Args) == 2 && os.Args[1] == ".") {
		operations.OpenFromCurrentDirectory()
		return
	} else if len(os.Args) == 2 {
		operations.OpenFromRepositoryName(os.Args[1])
		return
	}

	switch os.Args[1] {

	}

	fmt.Printf("neoman: '%s' is not a valid command. See 'nman --help'.\n", os.Args[1])
}
