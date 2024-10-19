package internal

import (
	"log"
	"os"
	"path"
)

func IsAlreadyRegistered(proj string) bool {
	registryDir, err := DocsRegistryDir()
	if err != nil {
		log.Fatalln(ErrAccessRegistryDir)
	}

	_, err = os.Stat(path.Join(registryDir, proj))
	return err == nil
}

func AddToRegistry(projName, projPath string) {}
