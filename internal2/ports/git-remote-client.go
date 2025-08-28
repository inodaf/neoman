package ports

import (
  "errors"
  "net/url"
)

var ErrGitRemoteNotFound = errors.New("could not find remote repo")

type GitRemoteClient interface {
  CloneURL(author, repo string) url.URL
  HasDocsDir(author, repo string) error
}
