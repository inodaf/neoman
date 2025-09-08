package operations

import (
	"fmt"

	"github.com/inodaf/neoman/internal3/management"
	"github.com/inodaf/neoman/packages/git"
)

func FetchDocs(owner, repo string) error {
	remote := git.NewGitHubClient()

	err := remote.IsDocsDirPresent(owner, repo)
	if err != nil {
		return fmt.Errorf("neoman: Could not locate 'docs/' from '%s/%s' on GitHub.\nMake sure you have reading rights", owner, repo)
	}

	fmt.Printf("neoman:	Fetching docs for '%s/%s' from GitHub...\n", owner, repo)

	return management.RegistryAddEntry(management.RegistryEntry{
		Scope:   management.RegistryTypeRemote,
		Owner:   owner,
		Project: repo,
	})
}
