package config

import (
	"os"
	"path"
)

// DocsRegistryDir returns the directory path that stores
// copies of docs. It creates the registry dir if that was not done before
// with the required sub-directories for "local" and "remote" documentation.
func DocsRegistryDir() (string, error) {
	appConfigDir, err := AppConfigDir()
	if err != nil {
		return "", err
	}

	docsRegistryDir := path.Join(appConfigDir, "registry")
	if _, err := os.Stat(docsRegistryDir); os.IsNotExist(err) {
		if err := os.Mkdir(docsRegistryDir, os.ModePerm); err != nil {
			return "", err
		}
		if err := os.Mkdir(path.Join(docsRegistryDir, "local"), os.ModePerm); err != nil {
			return "", err
		}
		if err := os.Mkdir(path.Join(docsRegistryDir, "remote"), os.ModePerm); err != nil {
			return "", err
		}
	}

	return docsRegistryDir, nil
}
