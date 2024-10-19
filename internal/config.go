package internal

import (
	"os"
	"path"
)

const AppName = "neoman"
const ShortAppName = "nman"

const PrimaryDocsDirName = "docs"
const AlternateDocsDirName = "manual"

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
		if err := os.Mkdir(appConfigDir, os.ModeDir); err != nil {
			return "", err
		}
	}

	return appConfigDir, nil
}

// DocsRegistryDir returns the directory path that stores
// copies of docs. It creates the dir if that was not done before.
func DocsRegistryDir() (string, error) {
	appConfigDir, err := AppConfigDir()
	if err != nil {
		return "", err
	}

	docsRegistryDir := path.Join(appConfigDir, "registry")
	if _, err := os.Stat(docsRegistryDir); os.IsNotExist(err) {
		if err := os.Mkdir(docsRegistryDir, os.ModeDir); err != nil {
			return "", err
		}
	}

	return docsRegistryDir, nil
}
