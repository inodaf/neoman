package internal

import (
	"fmt"
	"github.com/inodaf/neoman/internal/domains/ports"
	adapters2 "github.com/inodaf/neoman/internal/infra/adapters/driven/http/github"
	"github.com/inodaf/neoman/internal/infra/adapters/driven/persistence/registry"
)

func FetchDocs(owner, repo string) error {
	remote := adapters2.NewProviderGitHub()
	registry := adapters.FSProjectRegistry{}

	err := remote.DocsDirExists(owner, repo)
	if err != nil {
		return fmt.Errorf("neoman: Could not locate 'docs/' from '%s/%s' on GitHub.\nMake sure you have reading rights", owner, repo)
	}

	fmt.Printf("neoman:	Fetching docs for '%s/%s' from GitHub...\n", owner, repo)

	return registry.AddEntry(ports.RegistryEntry{
		Scope:   ports.RemoteScope,
		Owner:   owner,
		Project: repo,
	})
}
