package main

import (
	"fmt"
	"os"

	"github.com/inodaf/neoman/internal/commands"
	"github.com/inodaf/neoman/internal/daemon"
)

func main() {
	if err := daemon.IPC.Ping(); err != nil {
		fmt.Println("neoman: Could not connect to daemon")
		return
	}

	if len(os.Args) == 1 {
		commands.OpenFromWD()
		return
	} else if len(os.Args) == 2 {
		commands.OpenFromName(os.Args[1])
		return
	}

	switch os.Args[1] {

	}

	fmt.Printf("neoman: '%s' is not a valid command. See 'nman --help'.\n", os.Args[1])
}
