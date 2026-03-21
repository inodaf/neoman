package git

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

func (c *GitHubClient) ProviderName() string {
	return "github"
}

func (c *GitHubClient) CloneURL(author, repo string) string {
	var sshURL = url.URL{
		User: url.User("git"),
		Path: fmt.Sprintf("%s.git", repo),
		Host: fmt.Sprintf("%s:%s", "github.com", author),
	}

	return strings.Replace(sshURL.String(), "//", "", 1)
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
