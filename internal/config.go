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
	if _, err := os.Stat(appConfigDir); err != nil {
		if err := os.Mkdir(appConfigDir, os.ModePerm); err != nil {
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

	baseDir := path.Join(appConfigDir, "registry")
	if _, err := os.Stat(baseDir); err != nil {
		if err := os.Mkdir(baseDir, os.ModePerm); err != nil {
			return "", err
		}
	}

	return baseDir, nil
}
