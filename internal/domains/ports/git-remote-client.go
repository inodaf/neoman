package ports

import (
  "errors"
  "net/url"
)

var ErrGitRemoteNotFound = errors.New("neoman: Could not find repo")

type GitRemoteClient interface {
  CloneURL(owner, repo string) url.URL
  HasDocsDir(owner, repo string) error
}
