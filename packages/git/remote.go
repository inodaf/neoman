package git

import (
	"fmt"
	"net/http"
	"net/url"
)

type GitHubClient struct {
	baseRequest http.Request
}

func (p GitHubClient) IsDocsDirPresent(owner, repo string) error {
	p.baseRequest.URL.Path = fmt.Sprintf("repos/%s/%s/contents/docs", owner, repo)
	res, err := http.Get(p.baseRequest.URL.String())
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusNotFound {
		return ErrGitRemoteNotFound
	}

	return nil
}

func NewGitHubClient() *GitHubClient {
	h := make(http.Header, 2)
	h.Add("Accept", "application/vnd.github+json")
	h.Add("X-GitHub-Api-Version", "2022-11-28")

	return &GitHubClient{
		baseRequest: http.Request{
			Header: h,
			URL:    &url.URL{Scheme: "https", Host: "api.github.com"},
		},
	}
}
