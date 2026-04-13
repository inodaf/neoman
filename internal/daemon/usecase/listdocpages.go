package usecase

import (
	"fmt"

	"github.com/inodaf/neoman/internal/daemon/domain"
)

type ListDocPagesInput struct {
	Author     string
	Repository string
}

type ListDocPagesOutput struct {
	Pages []domain.DocsPage
}

func (u *UseCase) ListDocPages(input ListDocPagesInput) (ListDocPagesOutput, error) {
	if input.Author == "" {
		return ListDocPagesOutput{}, ErrListDocPagesAuthorRequired
	}

	if input.Repository == "" {
		return ListDocPagesOutput{}, ErrListDocPagesRepositoryRequired
	}

	exists, err := u.docsRepository.Exists(input.Author, input.Repository)
	if err != nil {
		return ListDocPagesOutput{}, err
	}

	if !exists {
		return ListDocPagesOutput{}, ErrListDocPagesDocsNotFound
	}

	pages, err := u.docsPageRepository.FindAll(input.Author, input.Repository)
	if err != nil {
		return ListDocPagesOutput{}, err
	}

	return ListDocPagesOutput{Pages: pages}, nil
}

var (
	ErrListDocPagesAuthorRequired     = fmt.Errorf("author is required")
	ErrListDocPagesRepositoryRequired = fmt.Errorf("repository is required")
	ErrListDocPagesDocsNotFound       = fmt.Errorf("docs not found")
)
