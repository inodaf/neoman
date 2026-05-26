package usecase

import "fmt"

type ListAuthorDocsInput struct {
	Author string
}

type ListAuthorDocsOutput struct {
	Documentations []string
}

func (u *UseCase) ListAuthorDocs(input ListAuthorDocsInput) (*ListAuthorDocsOutput, error) {
	if input.Author == "" {
		return nil, ErrListAuthorDocsAuthorRequired
	}

	docs, err := u.docsRepository.GetByAuthor(input.Author)
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, ErrListAuthorDocsNotFound
	}

	repositories := make([]string, len(docs))
	for i, doc := range docs {
		repositories[i] = doc.Repository
	}

	return &ListAuthorDocsOutput{Documentations: repositories}, nil
}

var (
	ErrListAuthorDocsAuthorRequired = fmt.Errorf("author is required")
	ErrListAuthorDocsNotFound       = fmt.Errorf("no documentation found for author")
)
