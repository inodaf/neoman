package git

type GitRemote interface {
	IsDocsDirPresent(author, repo string) error
	CloneURL(author, repo string) string
	Name() string
	WebPageURL(author, repo string) string
}
