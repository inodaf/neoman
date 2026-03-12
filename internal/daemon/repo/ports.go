package repo

import "github.com/inodaf/neoman/internal/daemon/domain"

type DocsRepository interface {
	Save(entry domain.RemoteDocs) error
	Exists(author, repository string) (bool, error)
	GetPath(entry domain.RemoteDocs) (string, error)
}
