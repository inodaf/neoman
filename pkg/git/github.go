package git

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
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

func (c *GitHubClient) HasDocs(repo Repo) (bool, error) {
	binPath, err := exec.LookPath("gh")
	if err != nil {
		return false, ErrGitHubCliNotInstalled
	}

	query := `query {
		repository(owner: "%s", name: "%s") { 
			name
			object(expression: "HEAD:docs") { __typename }
		}
  }`

	query = fmt.Sprintf(query, repo.Author, repo.Name)
	rawResult, err := exec.Command(binPath, "api", "graphql", "-f", fmt.Sprintf("query=%s", query)).Output()
	if err != nil {
		return false, ErrGitHubCliNotInstalled
	}

	var response struct {
		data struct {
			repository struct {
				object struct {
					Typename string `json:"__typename"`
				}
			}
		}
	}

	err = json.Unmarshal(rawResult, &response)
	if err != nil {
		return false, ErrGitHubUnableToParseResponse
	}

	return response.data.repository.object.Typename == "Tree", nil
}

func (c *GitHubClient) ListActiveRepos(author string, limit uint) ([]Repo, error) {
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

	response := make([]Repo, limit, 128)
	err = json.Unmarshal(rawResult, &response)
	if err != nil {
		return nil, ErrGitHubUnableToParseResponse
	}

	return response, nil
}

var (
	ErrGitHubCliNotInstalled       = errors.New("github cli was not found")
	ErrGitHubUnableToParseResponse = errors.New("unable to parse github api response")
)
