package usecases

import (
	"github.com/inodaf/neoman/internal2/ports"
)

type GetRemoteDocsUseCase struct {
	GitRemoteClient ports.GitRemoteClient
	ProjectRegistry ports.ProjectRegistry
	AuthorService   ports.AuthorService
}

func (u *GetRemoteDocsUseCase) Execute(authorName string, repoName string) error {
	author, err := u.AuthorService.FindOrCreate(authorName)
	if err != nil {
		return err
	}

	err = u.GitRemoteClient.HasDocsDir(authorName, repoName)
	if err != nil {
		return err
	}

	return u.ProjectRegistry.AddEntry(ports.RegistryEntry{
		Scope:   ports.RemoteScope,
		Owner:   author.Name,
		Project: repoName,
	})
}
