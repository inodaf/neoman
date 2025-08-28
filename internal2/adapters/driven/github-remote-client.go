package driven

import (
	"fmt"
	"github.com/inodaf/neoman/internal2/ports"
	"net/http"
	"net/url"
)

type GitHub struct {
	baseRequest http.Request
}

func (gh GitHub) CloneURL(author, repo string) url.URL {
	return url.URL{
		User: url.User("git"),
		Host: fmt.Sprintf("github.com:%s", author),
		Path: fmt.Sprintf("%s.git", repo),
	}
}

func (gh GitHub) HasDocsDir(author, repo string) error {
	gh.baseRequest.URL.Path = fmt.Sprintf("repos/%s/%s/contents/docs", author, repo)

	res, err := http.Get(gh.baseRequest.URL.String())
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusNotFound {
		return ports.ErrGitRemoteNotFound
	}

	return nil
}

func NewProviderGitHub() GitHub {
	header := make(http.Header, 2)

	header.Add("Accept", "application/vnd.github+json")
	header.Add("X-GitHub-Api-Version", "2022-11-28")

	return GitHub{
		baseRequest: http.Request{
			Header: header,
			URL:    &url.URL{Scheme: "https", Host: "api.github.com"},
		},
	}
}
