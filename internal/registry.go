package internal

import (
	"log"
	"os"
	"path"
)

func IsAlreadyRegistered(proj string, options ...string) bool {
	registryDir, err := DocsRegistryDir()
	if err != nil {
		log.Fatalln(ErrAccessRegistryDir)
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

func AddSymlinkToRegistry(proj, projPath string) error {
	if IsAlreadyRegistered(proj) {
		return ErrAlreadyRegistered
	}

	registryDir, err := DocsRegistryDir()
	if err != nil {
		return err
	}

	if err = os.Mkdir(path.Join(registryDir, "local", proj), os.ModePerm); err != nil {
		return ErrCreateLocalDocsDir
	}

	return os.Link(
		path.Join(projPath, PrimaryDocsDirName), // TODO: Handle when using AlternateDocsDirName.
		path.Join(registryDir, "local", proj, PrimaryDocsDirName),
	)
}
