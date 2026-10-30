package usecase

import "fmt"

func (u *UseCase) GetAll() (*GetAllOutput, error) {
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

	return &GetAllOutput{Docs: documentations}, nil
}

type GetAllOutput struct {
	Docs []string
}

var ErrListAllDocsEmpty = fmt.Errorf("no documentations found")
