package operations

import "github.com/inodaf/neoman/internal/management"

type availableDoc struct {
	Title        string
	RelativePath string
}

func ListAvailableDocs(author string, repo string) ([]availableDoc, error) {
	_, err := management.RegistryEntryDirPath(management.RegistryEntry{
		Scope:   management.RegistryTypeRemote,
		Project: repo,
		Owner:   author,
	})
	if err != nil {
		return nil, err
	}

	return make([]availableDoc, 0), nil
}
