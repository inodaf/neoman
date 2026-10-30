package usecase

import "fmt"

func (u *UseCase) GetFromAuthor(input GetFromAuthorInput) (*GetFromAuthorOutput, error) {
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

	return &GetFromAuthorOutput{Docs: repositories}, nil
}

type GetFromAuthorInput struct {
	Author string
}

type GetFromAuthorOutput struct {
	Docs []string
}

var (
	ErrListAuthorDocsAuthorRequired = fmt.Errorf("author is required")
	ErrListAuthorDocsNotFound       = fmt.Errorf("no documentation found for author")
)
