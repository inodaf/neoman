package config

import (
	"os"
	"path"
)

const AppName = "neoman"
const ShortAppName = "nman"

// AppSockAddr is the location of the Unix Domain Socket
// used for inter-process communication between
// Neoman's CLI and its daemon. Unix-only.
const AppSockAddr = "/tmp/nman.sock"
// const AppHostName = "neoman.local"
const AppHostName = "localhost:8092"
const AppWebAppPort = ":8092"
const AppDBFileName = "neoman.db"

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

// AppDataDir returns the directory path that is used
// for storing the database file.
// It creates the directory if that was not done before.
func AppDataDir() (string, error) {
	appConfigDir, err := AppConfigDir()
	if err != nil {
		return "", err
	}

	appDataDir := path.Join(appConfigDir, "data")
	if _, err := os.Stat(appDataDir); os.IsNotExist(err) {
		if err := os.Mkdir(appDataDir, os.ModePerm); err != nil {
			return "", err
		}
	}

	return appDataDir, nil
}