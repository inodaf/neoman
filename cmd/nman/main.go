package main

import (
	"fmt"
	"os"

	"github.com/inodaf/neoman/internal/commands"
	"github.com/inodaf/neoman/internal2/adapters/driven"
	"github.com/inodaf/neoman/internal2/adapters/driving"
)

var gitRemoteProvider = driven.NewProviderGitHub()
var projectRegistry = driven.FSProjectRegistry{GitRemoteClient: gitRemoteProvider}
var authorService = driven.IPCAuthorService{}
var cliPrompter = driving.CLIUserPrompter{}

func main() {
	err := driven.PingUnixSock()
	if err != nil {
		fmt.Println("neoman: Could not connect to daemon")
		return
	}

	controller := driving.DocsControllerCLI{
		gitRemoteProvider,
		projectRegistry,
		authorService,
		cliPrompter,
	}

	if len(os.Args) == 1 || (len(os.Args) == 2 && os.Args[1] == ".") {
		commands.OpenFromWD()
		return
	} else if len(os.Args) == 2 {
		err = controller.OpenFromRemote(os.Args[1])
		if err != nil {
			fmt.Printf("neoman: %s", err.Error())
		}

		return
	}

	switch os.Args[1] {

	}

	fmt.Printf("neoman: '%s' is not a valid command. See 'nman --help'.\n", os.Args[1])
}
