package main

import (
	"fmt"
	"os"

	"github.com/inodaf/neoman/internal/commands"
	"github.com/inodaf/neoman/internal/daemon"
)

func main() {
	daemon.IPC.Ping()

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
