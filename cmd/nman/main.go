package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/inodaf/neoman/internal/management"
	"github.com/inodaf/neoman/internal/operations"
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
	}

	if len(os.Args) == 2 && strings.Count(os.Args[1], "/") == 1 {
		operations.OpenFromName(os.Args[1])
		return
	}

	switch os.Args[1] {
		case "list":
			fmt.Println("neoman: 'list' command is not implemented yet. See 'nman --help' for more information.")
			return
		default:
			fmt.Printf("neoman: '%s' is not a valid command. See 'nman --help'.\n", os.Args[1])
			return
	}
}
