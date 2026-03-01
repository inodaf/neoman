package operations

import (
	"fmt"

	"github.com/inodaf/neoman/internal/management"
	"github.com/inodaf/neoman/pkg/git"
)

func GetDocs(owner, repo string) error {
	remote := git.NewGitHubClient()

	err := remote.IsDocsDirPresent(owner, repo)
	if err != nil {
		return fmt.Errorf("neoman: Could not locate 'docs/' from '%s/%s' on GitHub.\nMake sure you have reading rights", owner, repo)
	}

	fmt.Printf("neoman:	Fetching docs for '%s/%s' from GitHub...\n", owner, repo)

	err = management.RegistryAddEntry(management.RegistryEntry{
		Scope:   management.RegistryTypeRemote,
		Owner:   owner,
		Project: repo,
	})
	if err != nil {
		return err
	}

	// Sync docs to database after successful clone
	fmt.Printf("neoman:	Syncing docs to database for '%s/%s'...\n", owner, repo)

	db, err := management.NewSQLiteDatabase()
	if err != nil {
		return fmt.Errorf("neoman: could not open database for sync: %w", err)
	}
	defer db.Close()

	return Sync(owner, repo, db)
}
