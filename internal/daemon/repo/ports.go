package repo

import (
	"time"

	"github.com/inodaf/neoman/internal/daemon/domain"
)

type DocsRepository interface {
	Save(entry domain.RemoteDocs) error
	Exists(author, repository string) (bool, error)
	GetOne(author, repository string) (domain.RemoteDocs, error)

	StartIndexing(author, repository string) error
	StopIndexing(author, repository string) error
}

type DocsPageRepository interface {
	Save(page domain.DocsPage) error
	SaveMany(pages []domain.DocsPage) error
}

type SourceRegistry interface {
	Download(entry domain.RemoteDocs) error
	GetAllContents(entry domain.RemoteDocs) ([]RegistryContent, error)
}

type RegistryContent struct {
	Text           string
	RelPath        string
	LastModifiedAt time.Time
}
