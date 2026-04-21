package repo

import (
	"time"

	"github.com/inodaf/neoman/internal/daemon/domain"
)

type DocsRepository interface {
	Save(entry domain.RemoteDocs) error
	Exists(author, repository string) (bool, error)
	GetOne(author, repository string, source domain.RemoteSource) (domain.RemoteDocs, error)
	GetByAuthor(author string) ([]domain.RemoteDocs, error)
	GetAll() ([]domain.RemoteDocs, error)
}

type DocsPageRepository interface {
	Save(page domain.DocsPage) error
	SaveMany(pages []domain.DocsPage) error
	GetAll(author, repository string) ([]domain.DocsPage, error)
}

type ContentSource interface {
	Download(entry domain.RemoteDocs) error
	GetAll(entry domain.RemoteDocs) ([]RegistryContent, error)
}

type RegistryContent struct {
	Text           string
	RelPath        string
	LastModifiedAt time.Time
}

type JobRepository interface {
	Save(job *domain.Job) error
	GetPending(limit int) ([]domain.Job, error)
	GetAll() ([]domain.Job, error)
	Update(job *domain.Job) error
	Delete(jobID string) error
}
