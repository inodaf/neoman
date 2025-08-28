package driving

import (
	"errors"

	"github.com/inodaf/neoman/internal2/ports"
	"github.com/inodaf/neoman/internal2/usecases"
	"github.com/inodaf/neoman/internal2/valueobjects"
	"github.com/inodaf/neoman/packages/browser"
)

type DocsControllerCLI struct {
	GitRemoteClient ports.GitRemoteClient
	ProjectRegistry ports.ProjectRegistry
	AuthorService   ports.AuthorService
	UserPrompter    ports.UserPrompter
}

func (c *DocsControllerCLI) OpenFromRemote(proj string) error {
	authorAndRepo := valueobjects.AuthorAndRepo(proj)

	err := authorAndRepo.Validate()
	if err != nil {
		return err
	}

	authorName, repoName := authorAndRepo.Value()

	verifyAuthorUC := usecases.VerifyAuthorTrustUseCase{
		AuthorService: c.AuthorService,
		UserPrompter:  c.UserPrompter,
	}

	trusted, err := verifyAuthorUC.Execute(authorName)
	if err != nil {
		return err
	}

	if !trusted {
		return errors.New("author not trusted")
	}

	getDocsUC := usecases.GetRemoteDocsUseCase{
		GitRemoteClient: c.GitRemoteClient,
		ProjectRegistry: c.ProjectRegistry,
		AuthorService:   c.AuthorService,
	}

	err = getDocsUC.Execute(authorName, repoName)
	if err != nil {
		return err
	}

	err = browser.Open(proj)
	if err != nil {
		return err
	}

	return nil
}
