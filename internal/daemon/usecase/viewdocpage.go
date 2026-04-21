package usecase

import (
	"fmt"
	"strings"
)

type ViewDocPageInput struct {
	Author       string
	Repository   string
	RelativePath string
}

type ViewDocPageOutput struct {
	Content string
}

func (u *UseCase) ViewDocPage(input ViewDocPageInput) (*ViewDocPageOutput, error) {
	if input.Author == "" {
		return nil, ErrViewDocPageAuthorRequired
	}

	if input.Repository == "" {
		return nil, ErrViewDocPageRepositoryRequired
	}

	if input.RelativePath == "" {
		return nil, ErrViewDocPagePathRequired
	}

	// Normalize path: auto-append .md if no extension
	normalizedPath := input.RelativePath
	if !strings.HasSuffix(normalizedPath, ".md") && !strings.HasSuffix(normalizedPath, ".mdx") {
		normalizedPath = normalizedPath + ".md"
	}

	exists, err := u.docsRepository.Exists(input.Author, input.Repository)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrViewDocPageDocsNotFound
	}

	page, err := u.docsPageRepository.GetOne(input.Author, input.Repository, normalizedPath)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, ErrViewDocPageNotFound
		}
		return nil, err
	}

	return &ViewDocPageOutput{Content: page.Content}, nil
}

var (
	ErrViewDocPageAuthorRequired     = fmt.Errorf("author is required")
	ErrViewDocPageRepositoryRequired = fmt.Errorf("repository is required")
	ErrViewDocPagePathRequired       = fmt.Errorf("path is required")
	ErrViewDocPageDocsNotFound       = fmt.Errorf("documentation not found")
	ErrViewDocPageNotFound           = fmt.Errorf("page not found")
)
