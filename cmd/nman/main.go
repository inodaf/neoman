package main

import (
	"fmt"
	"log"
	"os"

	"github.com/inodaf/neoman/internal/commands"
	"github.com/inodaf/neoman/internal/daemon"
)

func main() {
	err := daemon.Ping()
	if err != nil {
		log.Fatal("neoman: could not reach daemon", err.Error())
	}

	if len(os.Args) == 1 {
		commands.OpenFromWD()
		return
	} else if len(os.Args) == 2 {
		commands.OpenFromName(os.Args[1])
		return
	}

	switch os.Args[1] {
	default:
		fmt.Printf("neoman: '%s' is not a valid command. See 'nman --help'.\n", os.Args[1])
		return
	}
}
