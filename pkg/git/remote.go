package git

import "time"

type GitRemote interface {
	Name() string
	CloneURL(author, repo string) string
	WebPageURL(author, repo string) string

	HasDocs(repo Repo) (bool, error)
	ListActiveRepos(author string, limit uint) ([]Repo, error)
}

type Repo struct {
	Author      string     `json:"author"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Visibility  visibility `json:"visibility"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type visibility string

const (
	VisibilityPublic  visibility = "public"
	VisibilityPrivate visibility = "private"
)
