package git

type GitRemote interface {
	IsDocsDirPresent(owner, repo string) error
}
