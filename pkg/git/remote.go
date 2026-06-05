package git

type GitRemote interface {
	Name() string
	CloneURL(author, repo string) string
	WebPageURL(author, repo string) string

	HasDocs(author string, repos []string) (map[string]bool, error)

	// CountActiveRepos(author string) (uint, error)
	// ListActiveRepos(author string, limit uint) ([]string, error)
}
