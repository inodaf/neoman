package git

import (
	"fmt"
	"net/http"
	"net/url"
)

func NewGitHubClient() *GitHubClient {
	h := make(http.Header, 2)

	h.Set("Accept", "application/vnd.github+json")
	h.Set("X-GitHub-Api-Version", "2022-11-28")

	return &GitHubClient{
		Request: http.Request{
			Header: h,
			URL:    &url.URL{Scheme: "https", Host: "api.github.com"},
		},
	}
}

type GitHubClient struct {
	http.Client
	http.Request
}

func (c *GitHubClient) Name() string {
	return "github"
}

func (c *GitHubClient) WebPageURL(author, repo string) string {
	return fmt.Sprintf("https://github.com/%s/%s", author, repo)
}

func (c *GitHubClient) CloneURL(author, repo string) string {
	return fmt.Sprintf("git@github.com:%s/%s.git", author, repo)
}

func (c *GitHubClient) IsDocsDirPresent(author, repo string) error {
	c.Request.URL.Path = fmt.Sprintf("repos/%s/%s/contents/docs", author, repo)
	res, err := c.Get(c.Request.URL.String())
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusNotFound {
		return ErrGitRemoteNotFound
	}

	return nil
}
