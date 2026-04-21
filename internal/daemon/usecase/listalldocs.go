package usecase

import "fmt"

type ListAllDocsOutput struct {
	Documentations []string
}

func (u *UseCase) ListAllDocs() (*ListAllDocsOutput, error) {
	docs, err := u.docsRepository.GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, ErrListAllDocsEmpty
	}

	documentations := make([]string, len(docs))
	for i, doc := range docs {
		documentations[i] = fmt.Sprintf("%s/%s", doc.Author, doc.Repository)
	}

	return &ListAllDocsOutput{Documentations: documentations}, nil
}

var (
	ErrListAllDocsEmpty = fmt.Errorf("No documentations found")
)
