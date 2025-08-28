package ports

import "github.com/inodaf/neoman/internal2/models"

type AuthorService interface {
	Trust(authorName string) error
	IsTrusted(authorName string) bool
	FindOrCreate(name string) (*models.Author, error)
}
