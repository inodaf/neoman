package usecase

import (
	"fmt"

	"github.com/inodaf/neoman/internal/daemon/domain"
)

type AddDocsInput struct {
	Author     string
	Repository string
}

func (u *UseCase) AddRemoteDocs(input AddDocsInput) error {
	if input.Author == "" {
		return ErrAddRemoteDocsAuthorRequired
	}

	if input.Repository == "" {
		return ErrAddRemoteDocsRepositoryRequired
	}

	exists, err := u.docsRepository.Exists(input.Author, input.Repository)
	if err != nil {
		return err
	}

	if exists {
		return ErrAddRemoteDocsAlreadyExist
	}

	err = u.gitRemoteClient.IsDocsDirPresent(input.Author, input.Repository)
	if err != nil {
		return ErrAddRemoteDocsDirNotFound
	}

	docs := domain.NewRemoteDocs(
		input.Author,
		input.Repository,
		domain.RemoteSource(u.gitRemoteClient.ProviderName()),
	)
	
	err = u.sourceRegistry.Download(*docs)
	if err != nil {
		return err
	}
	
	err = u.docsRepository.Save(*docs)
	if err != nil {
		return err
	}	

	return nil
}

var (
	ErrAddRemoteDocsAuthorRequired = fmt.Errorf("author is required")
	ErrAddRemoteDocsRepositoryRequired = fmt.Errorf("repository is required")
	ErrAddRemoteDocsAlreadyExist = fmt.Errorf("docs already exist")
	ErrAddRemoteDocsDirNotFound = fmt.Errorf("docs dir not found in repo")
)