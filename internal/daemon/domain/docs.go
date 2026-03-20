package domain

func NewRemoteDocs(author, repository string, source RemoteSource) *RemoteDocs {
	if author == "" || repository == "" {
		return nil
	}

	return &RemoteDocs{
		Source:     source,
		Author:     author,
		Repository: repository,
	}
}

type RemoteDocs struct {
	Source     RemoteSource
	Author     string
	Repository string
	Indexing   bool
}
