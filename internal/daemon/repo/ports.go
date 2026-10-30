package repo

import (
	"time"

	"github.com/inodaf/neoman/internal/daemon/domain"
)

type DocsRepository interface {
	Save(entry domain.Docs) error
	Exists(author, repository string) (bool, error)
	GetOne(author, repository string, source domain.ContentSource) (domain.Docs, error)
	GetByAuthor(author string) ([]domain.Docs, error)
	GetAll() ([]domain.Docs, error)
}

type DocsPageRepository interface {
	Save(page domain.Page) error
	SaveMany(pages []domain.Page) error
	GetAll(author, repository string) ([]domain.Page, error)
	GetOne(author, repository, relativePath string) (*domain.Page, error)
}

type ContentSource interface {
	Download(entry domain.Docs) error
	GetAll(entry domain.Docs) ([]RegistryContent, error)
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
