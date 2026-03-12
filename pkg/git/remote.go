package git

type GitRemote interface {
	IsDocsDirPresent(author, repo string) error
	ProviderName() string
}
