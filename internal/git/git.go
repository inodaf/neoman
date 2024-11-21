package git

import (
	"errors"
	"os/exec"
)

var ErrGitNotInstalled = errors.New("neoman: Git is not installed or was not found")
var ErrGitCloneForbidden = errors.New("neoman: Could not perform clone - forbidden")
var ErrGitRemoteNotFound = errors.New("neoman: Could not find repo")

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

func Clone(uri string) error {
	binPath, err := exec.LookPath("git")
	if err != nil {
		return ErrGitNotInstalled
	}

	_, err = exec.Command(binPath, "clone", uri, "--depth", "1").Output()
	return err
}
