package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/inodaf/neoman/internal/app/command"
	"github.com/inodaf/neoman/internal/management"
	"github.com/inodaf/neoman/internal/operations"
)

func main() {
	cmd := command.New()

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
		cmd.AddOrOpen(context.TODO(), os.Args[1])
		return
	}

	switch os.Args[1] {
	case "list":
		if len(os.Args) != 3 || strings.Count(os.Args[2], "/") != 1 {
			fmt.Println("neoman: Usage: nman list <author/repo>")
			return
		}
		cmd.ListPages(context.TODO(), os.Args[2])
		return
	case "view":
		handleViewCommand()
		return
	default:
		fmt.Printf("neoman: '%s' is not a valid command. See 'nman --help'.\n", os.Args[1])
		return
	}
}

// handleViewCommand processes the view command
// Usage: nman view <project> <document-path>
func handleViewCommand() {
	if len(os.Args) != 4 {
		fmt.Println("neoman: Usage: nman view <project> <document-path>")
		return
	}

	project := os.Args[2]
	documentPath := os.Args[3]

	// Call ViewDocument
	output, err := operations.ViewDocument(project, documentPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "neoman: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(output)
}
