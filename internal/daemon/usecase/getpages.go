package usecase

import (
	"fmt"

	"github.com/inodaf/neoman/internal/daemon/domain"
)

func (u *UseCase) GetPages(input GetPagesInput) (*GetPagesOutput, error) {
	if input.Author == "" {
		return nil, ErrListDocPagesAuthorRequired
	}

	if input.Repository == "" {
		return nil, ErrListDocPagesRepositoryRequired
	}

	exists, err := u.docsRepository.Exists(input.Author, input.Repository)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrListDocPagesDocsNotFound
	}

	pages, err := u.docsPageRepository.GetAll(input.Author, input.Repository)
	if err != nil {
		return nil, err
	}

	return &GetPagesOutput{Pages: pages}, nil
}

type GetPagesInput struct {
	Author     string
	Repository string
}

type GetPagesOutput struct {
	Pages []domain.Page
}

var (
	ErrListDocPagesAuthorRequired     = fmt.Errorf("author is required")
	ErrListDocPagesRepositoryRequired = fmt.Errorf("repository is required")
	ErrListDocPagesDocsNotFound       = fmt.Errorf("docs not found")
)
