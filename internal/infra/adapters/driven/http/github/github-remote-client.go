package adapters

import (
	"fmt"
	"github.com/inodaf/neoman/internal/domains/ports"
	"net/http"
	"net/url"
)

type GitHub struct {
	baseRequest http.Request
}

func (p GitHub) CloneURL(owner, repo string) url.URL {
	return url.URL{
		User: url.User("git"),
		Host: fmt.Sprintf("github.com:%s", owner),
		Path: fmt.Sprintf("%s.git", repo),
	}
}

func (p GitHub) DocsDirExists(owner, repo string) error {
	p.baseRequest.URL.Path = fmt.Sprintf("repos/%s/%s/contents/docs", owner, repo)
	res, err := http.Get(p.baseRequest.URL.String())
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusNotFound {
		return ports.ErrGitRemoteNotFound
	}

	return nil
}

func NewProviderGitHub() *GitHub {
	h := make(http.Header, 2)
	h.Add("Accept", "application/vnd.github+json")
	h.Add("X-GitHub-Api-Version", "2022-11-28")

	return &GitHub{
		baseRequest: http.Request{
			Header: h,
			URL:    &url.URL{Scheme: "https", Host: "api.github.com"},
		},
	}
}
