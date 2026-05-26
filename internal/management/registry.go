package management

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path"

	"github.com/inodaf/neoman/pkg/config"
	"github.com/inodaf/neoman/pkg/git"
)

type registryType int

const (
	RegistryTypeLocal registryType = iota
	RegistryTypeRemote
)

type RegistryEntry struct {
	Scope       registryType
	Owner       string
	Project     string
	ProjectPath string
}

func RegistryHasEntry(entry RegistryEntry) bool {
	registryDir, err := config.DocsRegistryDir()
	if err != nil {
		log.Fatalln("Could not access documentation registry.")
	}

	var entryPath string

	switch entry.Scope {
	case RegistryTypeLocal:
		entryPath = path.Join(registryDir, "local", entry.Project)
	case RegistryTypeRemote:
		entryPath = path.Join(registryDir, "remote", entry.Owner, entry.Project)
	}

	if len(entryPath) == 0 {
		log.Fatalln("Could not resolve Entry's path.")
	}

	_, err = os.Stat(entryPath)
	return err == nil
}

func RegistryEntryDirPath(entry RegistryEntry) (string, error) {
	registryDir, err := config.DocsRegistryDir()
	if err != nil {
		return "", err
	}

	switch entry.Scope {
	case RegistryTypeLocal:
		return path.Join(registryDir, "local", entry.Project), nil
	case RegistryTypeRemote:
		return path.Join(registryDir, "remote", entry.Owner, entry.Project), nil
	default:
		return "", errors.New("invalid registry type")
	}
}

func RegistryAddEntry(entry RegistryEntry) error {
	registryDir, err := config.DocsRegistryDir()
	if err != nil {
		return err
	}

	switch entry.Scope {
	case RegistryTypeLocal:
		return addLocalEntry(entry, registryDir)
	case RegistryTypeRemote:
		return addRemoteEntry(entry, registryDir)
	default:
		return nil
	}
}

func addLocalEntry(entry RegistryEntry, registryDir string) error {
	if err := os.Mkdir(path.Join(registryDir, "local", entry.Project), os.ModePerm); err != nil {
		return errors.New("could not create 'local' registry directory for project")
	}

	return os.Symlink(
		path.Join(entry.ProjectPath, config.PrimaryDocsDirName),
		path.Join(registryDir, "local", entry.Project, config.PrimaryDocsDirName),
	)
}

func addRemoteEntry(input RegistryEntry, registryDir string) error {
	ownerDir := path.Join(registryDir, "remote", input.Owner)
	if _, err := os.Stat(ownerDir); errors.Is(err, fs.ErrNotExist) {
		if err := os.Mkdir(ownerDir, os.ModePerm); err != nil {
			return err
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return errors.New("could not get current working directory")
	}

	err = os.Chdir(ownerDir)
	if err != nil {
		return err
	}

	defer os.Chdir(wd)

	err = git.Clone(input.Owner, input.Project, git.NewGitHubClient())
	if err != nil {
		return errors.New("could not copy documentation from remote")
	}

	err = os.Chdir(path.Join(ownerDir, input.Project))
	if err != nil {
		return errors.New("could not access documentation directory")
	}

	err = git.SparseCheckout()
	if err != nil {
		return errors.New("could not sync documentation from remote")
	}

	return nil
}
