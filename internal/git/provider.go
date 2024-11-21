package git

type RemoteProvider interface {
	GitURL(owner, repo string) string
	DocsDirExists(owner, repo string) error
}
