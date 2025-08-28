package internal

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path"
	
	"github.com/inodaf/neoman/internal2/adapters/driven"
)

func IsAlreadyRegistered(proj string, options ...string) bool {
	registryDir, err := DocsRegistryDir()
	if err != nil {
		log.Fatalln(driven.ErrAccessRegistryDir)
	}

	registryScope := "local"
	if len(options) > 0 {
		if options[0] == "local" || options[0] == "remote" {
			registryScope = options[0]
		}
	}

	_, err = os.Stat(path.Join(registryDir, registryScope, proj))
	return err == nil
}

func AddLocalEntryToRegistry(proj, projPath string) error {
	if IsAlreadyRegistered(proj) {
		return driven.ErrAlreadyRegistered
	}

	registryDir, err := DocsRegistryDir()
	if err != nil {
		return err
	}

	if err = os.Mkdir(path.Join(registryDir, "local", proj), os.ModePerm); err != nil {
		return driven.ErrCreateLocalDocsDir
	}

	return os.Symlink(
		path.Join(projPath, PrimaryDocsDirName),
		path.Join(registryDir, "local", proj, PrimaryDocsDirName),
	)
}

func AddRemoteEntryToRegistry(owner string) (string, error) {
	dir, err := DocsRegistryDir()
	if err != nil {
		return "", err
	}

	entry := path.Join(dir, "remote", owner)
	if _, err := os.Stat(entry); errors.Is(err, fs.ErrNotExist) {
		err := os.Mkdir(entry, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return entry, nil
}
