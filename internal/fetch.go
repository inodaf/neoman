package internal

import (
	"fmt"
	"log"
	"os"

	"github.com/inodaf/neoman/internal/git"
)

func FetchDocs(owner, repo string) error {
	remote := git.NewProviderGitHub()

	err := remote.DocsDirExists(owner, repo)
	if err != nil {
		return fmt.Errorf("neoman: Could not locate 'docs/' from '%s/%s' on GitHub.\nMake sure you have reading rights", owner, repo)
	}

	fmt.Printf("neoman:	Fetching docs for '%s/%s' from GitHub...\n", owner, repo)
	ownerRegistryDir, err := AddRemoteEntryToRegistry(owner)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(ErrGetWd)
	}

	err = os.Chdir(ownerRegistryDir)
	if err != nil {
		return err
	}
	defer os.Chdir(wd)

	err = git.Clone(remote.GitURL(owner, repo))
	if err != nil {
		return err
	}

	return nil
}
