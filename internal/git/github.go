package git

import (
	"fmt"
	"net/http"
	"net/url"
)

type GitHub struct {
	baseRequest http.Request
}

func (p GitHub) GitURL(owner, repo string) string {
	return fmt.Sprintf("git@github.com:%s/%s.git", owner, repo)
}

func (p GitHub) DocsDirExists(owner, repo string) error {
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

func NewProviderGitHub() *GitHub {
	h := make(http.Header)
	h.Add("Accept", "application/vnd.github+json")
	h.Add("X-GitHub-Api-Version", "2022-11-28")

	return &GitHub{
		baseRequest: http.Request{
			Header: h,
			URL:    &url.URL{Scheme: "https", Host: "api.github.com"},
		},
	}
}
