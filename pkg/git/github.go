package git

import (
	"fmt"
	"net/http"
	"net/url"
)

type GitHubClient struct {
	http.Client
	http.Request
}

func (client *GitHubClient) IsDocsDirPresent(owner, repo string) error {
	client.Request.URL.Path = fmt.Sprintf("repos/%s/%s/contents/docs", owner, repo)
	res, err := client.Get(client.Request.URL.String())
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
		Request: http.Request{
			Header: h,
			URL:    &url.URL{Scheme: "https", Host: "api.github.com"},
		},
	}
}
