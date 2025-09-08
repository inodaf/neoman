package git

import (
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

var (
	ErrGitNotInstalled   = errors.New("git not installed or was not found")
	ErrGitCloneForbidden = errors.New("could not perform clone - forbidden")
	ErrGitRemoteNotFound = errors.New("could not find repo")
)

type GitRemoteProvider string

const (
	GitRemoteProviderGitHub    GitRemoteProvider = "github.com"
	GitRemoteProviderGitLab    GitRemoteProvider = "gitlab.com"
	GitRemoteProviderBitbucket GitRemoteProvider = "bitbucket.org"
)

// IsRepository checks if the current working directory
// has a valid Git repository by checking the "git status" command
// output. It returns [ErrGitNotInstalled] if "git" is not located in PATH.
func IsRepository() (bool, error) {
	binPath, err := exec.LookPath("git")
	if err != nil {
		return false, ErrGitNotInstalled
	}

	_, err = exec.Command(binPath, "status").Output()
	return err == nil, nil
}

// Clone fetches contents of a repository from a remote source.
// To preserve disk space and improve download speed only the
// last commit is downloaded. It returns [ErrGitNotInstalled] if "git" is not
// located in PATH.
func Clone(author, repo string, provider GitRemoteProvider) error {
	binPath, err := exec.LookPath("git")
	if err != nil {
		return ErrGitNotInstalled
	}

	var sshURL = url.URL{
		User: url.User("git"),
		Path: fmt.Sprintf("%s.git", repo),
		Host: fmt.Sprintf("%s:%s", provider, author),
	}

	cloneURL := strings.Replace(sshURL.String(), "//", "", 1)
	_, err = exec.Command(binPath, "clone", cloneURL, "--depth", "1").Output()
	return err
}
