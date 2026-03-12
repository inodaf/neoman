package usecase

import (
	"fmt"

	"github.com/inodaf/neoman/internal/daemon/domain"
)

type AddDocsInput struct {
	Author     string
	Repository string
}

type AddDocsError error

var ErrAuthorRequired AddDocsError = fmt.Errorf("author is required")
var ErrRepositoryRequired AddDocsError = fmt.Errorf("repository is required")
var ErrDocsAlreadyExist AddDocsError = fmt.Errorf("docs already exist")
var ErrDocsDirNotFound AddDocsError = fmt.Errorf("docs dir not found in repo")

func (u *UseCase) AddRemoteDocs(input *AddDocsInput) AddDocsError {
	if input.Author == "" {
		return ErrAuthorRequired
	}

	if input.Repository == "" {
		return ErrRepositoryRequired
	}

	exists, err := u.docsRepository.Exists(input.Author, input.Repository)
	if err != nil {
		return err
	}

	if exists {
		return ErrDocsAlreadyExist
	}

	err = u.gitRemoteClient.IsDocsDirPresent(input.Author, input.Repository)
	if err != nil {
		return ErrDocsDirNotFound
	}

	docs := domain.NewRemoteDocs(
		input.Author,
		input.Repository,
		domain.RemoteSource(u.gitRemoteClient.ProviderName()),
	)
	
	err = u.docsRepository.Save(*docs)
	if err != nil {
		return err
	}	

	return nil
}
