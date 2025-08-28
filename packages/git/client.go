package git

import (
	"errors"
	"net/url"
	"os/exec"
	"strings"
)

var ErrGitNotInstalled = errors.New("git is not installed or was not found")
var ErrGitCloneForbidden = errors.New("could not perform clone - forbidden")
var ErrGitRemoteNotFound = errors.New("could not find repo")

// IsRepository checks if the current working directory
// is a valid Git repository by checking the "git status" command
// output. It [Exit(1)] if the "git" binary was not located in PATH.
func IsRepository() (bool, error) {
	binPath, err := exec.LookPath("git")
	if err != nil {
		return false, ErrGitNotInstalled
	}

	_, err = exec.Command(binPath, "status").Output()
	return err == nil, nil
}

// Clone fetches contents of a repository from a remote source.
// To preserve disk space and improve download speed only the last commit is downloaded.
func Clone(remoteURL url.URL) error {
	binPath, err := exec.LookPath("git")
	if err != nil {
		return ErrGitNotInstalled
	}

	cloneURL := strings.Replace(remoteURL.String(), "//", "", 1)
	_, err = exec.Command(binPath, "clone", cloneURL, "--depth", "1").Output()
	return err
}
