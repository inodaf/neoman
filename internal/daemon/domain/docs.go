package domain

import (
	"fmt"
	"time"
)

func NewRemoteDocs(author, repository string, source RemoteSource) *RemoteDocs {
	if author == "" || repository == "" {
		return nil
	}

	now := time.Now()
	return &RemoteDocs{
		Source:     source,
		Author:     author,
		Repository: repository,
		Indexing:   false,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}
}

type RemoteDocs struct {
	Source     RemoteSource
	Author     string
	Repository string
	Indexing   bool
	LastSyncAt *time.Time
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

// SetIndexing sets the indexing flag
func (r *RemoteDocs) SetIndexing(value bool) {
	r.Indexing = value
}

// StartIndexing marks this RemoteDocs as currently indexing
func (r *RemoteDocs) StartIndexing() {
	r.Indexing = true
}

// StopIndexing marks this RemoteDocs as done indexing and updates last sync time
func (r *RemoteDocs) StopIndexing() {
	r.Indexing = false
	now := time.Now()
	r.LastSyncAt = &now
}

// GetID returns composite ID for this RemoteDocs
func (r *RemoteDocs) GetID() string {
	return fmt.Sprintf("%s-%s", r.Author, r.Repository)
}
