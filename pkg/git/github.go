package git

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
)

func NewGitHubClient() GitRemote {
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

func (c *GitHubClient) HasDocs(author string, repos []string) (map[string]bool, error) {
	binPath, err := exec.LookPath("gh")
	if err != nil {
		return nil, ErrGitHubCliNotInstalled
	}

	query := `query {
		%s
  }
  fragment DocsFragment on Repository {
    name
    object(expression: "HEAD:docs") {
      __typename
    }
  }`

	var operations strings.Builder
	for i, repo := range repos {
		operation := `
    repo%d: repository(owner: "%s", name: "%s") { ...DocsFragment }
		`
		fmt.Fprintf(&operations, operation, i, author, repo)
	}

	query = fmt.Sprintf(query, operations.String())
	rawResult, err := exec.Command(binPath, "api", "graphql", "-f", fmt.Sprintf("query=%s", query)).Output()

	var response struct {
		Data map[string]struct {
			Name   string `json:"name"`
			Object struct {
				Typename string `json:"__typename"`
			} `json:"object"`
		} `json:"data"`
	}

	err = json.Unmarshal(rawResult, &response)
	if err != nil {
		return nil, ErrGitHubUnableToParseResponse
	}

	availability := make(map[string]bool, len(repos))
	for _, repo := range response.Data {
		availability[repo.Name] = repo.Object.Typename != "null"
	}

	return availability, nil
}

func (c *GitHubClient) ListActiveRepos(author string, limit uint) ([]string, error) {
	binPath, err := exec.LookPath("gh")
	if err != nil {
		return nil, ErrGitHubCliNotInstalled
	}

	args := []string{
		"repo",
		"list",
		author,
		"--no-archived",
		"--source",
		"--json",
		"name,description,visibility,updatedAt",
	}

	if limit > 0 {
		args = append(args, "--limit", fmt.Sprintf("%d", limit))
	}

	rawResult, err := exec.Command(binPath, args...).Output()
	if err != nil {
		return nil, ErrGitHubUnableToParseResponse
	}
	
	var response []Repo
	err = json.Unmarshal(rawResult, &response)
	//todo: resume from here
}

var ErrGitHubCliNotInstalled = errors.New("GitHub CLI was not found")
var ErrGitHubUnableToParseResponse = errors.New("Unable to parse GitHub API response")
