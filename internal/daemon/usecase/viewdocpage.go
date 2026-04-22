package usecase

import (
	"fmt"
	"strings"

	"github.com/inodaf/neoman/internal/daemon/domain"
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

	// Check if documentation exists
	exists, err := u.docsRepository.Exists(input.Author, input.Repository)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrViewDocPageDocsNotFound
	}

	// Normalize path: if no extension, try .md first, then .mdx
	pathsToTry := []string{input.RelativePath}
	if !strings.HasSuffix(input.RelativePath, ".md") && !strings.HasSuffix(input.RelativePath, ".mdx") {
		pathsToTry = []string{
			input.RelativePath + ".md",
			input.RelativePath + ".mdx",
		}
	}

	// Try each path until we find the page
	var page *domain.DocsPage
	for _, path := range pathsToTry {
		page, err = u.docsPageRepository.GetOne(input.Author, input.Repository, path)
		if err == nil {
			break
		}
		if !strings.Contains(err.Error(), "not found") {
			return nil, err
		}
	}

	if page == nil {
		return nil, ErrViewDocPageNotFound
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
