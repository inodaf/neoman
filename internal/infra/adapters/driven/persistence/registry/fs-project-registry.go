package adapters

import (
	"errors"
	"github.com/inodaf/neoman/internal/domains/ports"
	"github.com/inodaf/neoman/packages/config"
	"github.com/inodaf/neoman/packages/git"
	"io/fs"

	"log"
	"os"
	"path"
)

var (
	ErrGetWd              = errors.New("neoman: Could not get current working directory")
	ErrAccessRegistryDir  = errors.New("neoman: Could not access documentation registry")
	ErrAlreadyRegistered  = errors.New("neoman: A project of same name is already registered")
	ErrCreateLocalDocsDir = errors.New("neoman: Could not create 'local' registry directory for project")
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

	switch {
	case entry.Scope == ports.LocalScope:
		return i.addLocalEntry(entry, registryDir)
	case entry.Scope == ports.RemoteScope:
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
		return os.Mkdir(ownerDir, os.ModePerm)
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
