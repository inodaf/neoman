package main

import (
	"fmt"
	"log"
	"os"

	"github.com/inodaf/neoman/internal/commands"
)

func main() {
	if len(os.Args) == 1 {
		commands.OpenFromWD()
		return
	} else if len(os.Args) == 2 {
		log.Printf("try docs/ from '%s'", os.Args[1])
		return
	}

	switch os.Args[1] {
	default:
		fmt.Printf("neoman: '%s' is not a valid command. See 'nman --help'.\n", os.Args[1])
		return
	}
}
