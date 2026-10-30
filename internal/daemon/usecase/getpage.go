package usecase

import (
	"fmt"
	"strings"

	"github.com/inodaf/neoman/internal/daemon/domain"
)

func (u *UseCase) GetPage(input GetPageInput) (*GetPageOutput, error) {
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
	var page *domain.Page
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

	return &GetPageOutput{Content: page.Content}, nil
}

type GetPageInput struct {
	Author       string
	Repository   string
	RelativePath string
}

type GetPageOutput struct {
	Content string
}

var (
	ErrViewDocPageAuthorRequired     = fmt.Errorf("author is required")
	ErrViewDocPageRepositoryRequired = fmt.Errorf("repository is required")
	ErrViewDocPagePathRequired       = fmt.Errorf("path is required")
	ErrViewDocPageDocsNotFound       = fmt.Errorf("documentation not found")
	ErrViewDocPageNotFound           = fmt.Errorf("page not found")
)
