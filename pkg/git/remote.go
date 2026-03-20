package git

type GitRemote interface {
	IsDocsDirPresent(author, repo string) error
	CloneURL(author, repo string) string
	ProviderName() string
}
