package domain

import (
	"fmt"
	"time"
)

func NewDocs(author, repository string, source ContentSource) *Docs {
	if author == "" || repository == "" {
		return nil
	}

	now := time.Now()
	return &Docs{
		Source:     source,
		Author:     author,
		Repository: repository,
		Indexing:   false,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}
}

type Docs struct {
	Source     ContentSource
	Author     string
	Repository string
	Indexing   bool
	LastSyncAt *time.Time
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

// SetIndexing sets the indexing flag
func (r *Docs) SetIndexing(value bool) {
	r.Indexing = value
}

// StartIndexing marks this RemoteDocs as currently indexing
func (r *Docs) StartIndexing() {
	r.Indexing = true
}

// StopIndexing marks this RemoteDocs as done indexing and updates last sync time
func (r *Docs) StopIndexing() {
	r.Indexing = false
	now := time.Now()
	r.LastSyncAt = &now
}

// GetID returns composite ID for this RemoteDocs
func (r *Docs) GetID() string {
	return fmt.Sprintf("%s-%s", r.Author, r.Repository)
}
