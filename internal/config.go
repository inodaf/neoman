package internal

import (
	"os"
	"path"
)

const AppName = "neoman"
const ShortAppName = "nman"

// AppSockAddr is the location of the UDS
// used for inter-process communication between
// Neoman's CLI and its daemon. Unix-only.
const AppSockAddr = "/tmp/nman.sock"
const AppHostName = "neoman.local"

// PrimaryDocsDirName dictates Neoman's convention of having
// a `docs/` directory at the root of a project.
const PrimaryDocsDirName = "docs"

// AppConfigDir returns the directory path that used
// for storing internal app support files.
// It creates the directory if that was not done before.
func AppConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appConfigDir := path.Join(configDir, AppName)
	if _, err := os.Stat(appConfigDir); os.IsNotExist(err) {
		if err := os.Mkdir(appConfigDir, os.ModePerm); err != nil {
			return "", err
		}
	}

	return appConfigDir, nil
}

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
