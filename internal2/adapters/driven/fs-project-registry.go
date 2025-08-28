package driven

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path"

	"github.com/inodaf/neoman/internal2/ports"
	"github.com/inodaf/neoman/packages/config"
	"github.com/inodaf/neoman/packages/git"
)

var (
	ErrGetWd              = errors.New("could not get current working directory")
	ErrAccessRegistryDir  = errors.New("could not access documentation registry")
	ErrAlreadyRegistered  = errors.New("a project of same name is already registered")
	ErrCreateLocalDocsDir = errors.New("could not create 'local' registry directory for project")
)

type FSProjectRegistry struct {
	GitRemoteClient ports.GitRemoteClient
}

func (i FSProjectRegistry) HasEntry(entry ports.RegistryEntry) bool {
	registryDir, err := config.DocsRegistryDir()
	if err != nil {
		log.Fatalln(ErrAccessRegistryDir)
	}

	var entryPath string
	switch entry.Scope {
	case ports.LocalScope:
		entryPath = path.Join(registryDir, "local", entry.Project)
	case ports.RemoteScope:
		entryPath = path.Join(registryDir, "remote", entry.Owner, entry.Project)
	}

	if len(entryPath) == 0 {
		log.Fatalln("Could not resolve Entry's path.")
	}

	_, err = os.Stat(entryPath)
	return err == nil
}

func (i FSProjectRegistry) AddEntry(entry ports.RegistryEntry) error {
	if i.HasEntry(entry) {
		return ErrAlreadyRegistered
	}

	registryDir, err := config.DocsRegistryDir()
	if err != nil {
		return err
	}

	switch entry.Scope {
	case ports.LocalScope:
		return i.addLocalEntry(entry, registryDir)
	case ports.RemoteScope:
		return i.addRemoteEntry(entry, registryDir)
	default:
		return nil
	}
}

func (i FSProjectRegistry) addLocalEntry(entry ports.RegistryEntry, registryDir string) error {
	if err := os.Mkdir(path.Join(registryDir, "local", entry.Project), os.ModePerm); err != nil {
		return ErrCreateLocalDocsDir
	}

	return os.Symlink(
		path.Join(entry.ProjectPath, config.PrimaryDocsDirName),
		path.Join(registryDir, "local", entry.Project, config.PrimaryDocsDirName),
	)
}

func (i FSProjectRegistry) addRemoteEntry(input ports.RegistryEntry, registryDir string) error {
	ownerDir := path.Join(registryDir, "remote", input.Owner)
	if _, err := os.Stat(ownerDir); errors.Is(err, fs.ErrNotExist) {
		if err := os.Mkdir(ownerDir, os.ModePerm); err != nil {
			return err
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return ErrGetWd
	}

	err = os.Chdir(ownerDir)
	if err != nil {
		return err
	}
	defer os.Chdir(wd)

	return git.Clone(i.GitRemoteClient.CloneURL(input.Owner, input.Project))
}
