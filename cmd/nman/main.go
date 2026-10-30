package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/inodaf/neoman/internal/app/command"
	"github.com/inodaf/neoman/internal/operations"
)

func main() {
	cmd := command.New()

	err := cmd.SocketClientPing()
	if err != nil {
		fmt.Println("neoman: Could not connect to daemon")
		return
	}

	// Handle: nman or nman .
	if len(os.Args) == 1 || (len(os.Args) == 2 && os.Args[1] == ".") {
		operations.OpenFromCurrentDirectory()
		return
	}

	switch os.Args[1] {
	case "add":
		// Handle: nman add <author/repo>
		if len(os.Args) != 3 || strings.Count(os.Args[1], "/") != 1 {
			fmt.Println("neoman: Usage: nman add <author/repo>")
			return
		}

		cmd.AddOnly(context.TODO(), os.Args[1])
		return
	case "list":
		if len(os.Args) == 2 {
			// nman list → list all docs
			cmd.ListAllDocs(context.TODO())
			return
		}
		if len(os.Args) != 3 {
			fmt.Println("neoman: Usage: nman list, nman list <author>, or nman list <author/repo>")
			return
		}
		arg := os.Args[2]
		if strings.Count(arg, "/") == 1 {
			cmd.ListPages(context.TODO(), arg)
		} else if strings.Count(arg, "/") == 0 {
			cmd.ListAuthorDocs(context.TODO(), arg)
		} else {
			fmt.Println("neoman: Usage: nman list, nman list <author>, or nman list <author/repo>")
		}
		return
	case "view":
		if len(os.Args) != 4 {
			fmt.Println("neoman: Usage: nman view <author/repo> <path>")
			os.Exit(1)
			return
		}
		cmd.ViewPage(context.TODO(), os.Args[2], os.Args[3])
		return
	case "skills":
		cmd.Skills()
		return
	default:
		fmt.Printf("neoman: '%s' is not a valid command. See 'nman --help'.\n", os.Args[1])
		return
	}
}
