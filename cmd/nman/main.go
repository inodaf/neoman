package main

import (
	"context"
	"fmt"
	"io"
	"net/url"
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
		handleListCommand()
		return
	case "view":
		handleViewCommand()
		return
	default:
		fmt.Printf("neoman: '%s' is not a valid command. See 'nman --help'.\n", os.Args[1])
		return
	}
}

// handleListCommand processes the list command with three scenarios:
// 1. nman list - List all projects
// 2. nman list org - List projects in organization
// 3. nman list owner/repo - List documents in project
func handleListCommand() {
	var endpoint string

	if len(os.Args) == 2 {
		// Scenario 1: nman list - List all projects
		endpoint = "/list"
	} else if len(os.Args) == 3 {
		arg := os.Args[2]
		slashCount := strings.Count(arg, "/")

		if slashCount == 0 {
			// Scenario 2: nman list org - List projects in organization
			endpoint = "/list/" + arg
		} else if slashCount == 1 {
			// Scenario 3: nman list owner/repo - List documents in project
			parts := strings.Split(arg, "/")
			endpoint = "/list/" + parts[0] + "/" + parts[1]
		} else {
			fmt.Println("neoman: Invalid argument. Use 'nman list', 'nman list org', or 'nman list org/repo'")
			return
		}
	} else {
		fmt.Println("neoman: Too many arguments. Use 'nman list', 'nman list org', or 'nman list org/repo'")
		return
	}

	// Make HTTP GET request to daemon
	resource := url.URL{Host: "unix", Scheme: "http", Path: endpoint}
	resp, err := management.UnixSockClient.Get(resource.String())
	if err != nil {
		fmt.Println("neoman: Could not connect to daemon")
		return
	}
	defer resp.Body.Close()

	// Read and display response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("neoman: Could not read daemon response")
		return
	}

	fmt.Print(string(body))
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
