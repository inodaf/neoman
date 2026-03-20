package git

import (
	"errors"
	"os/exec"

	"github.com/inodaf/neoman/pkg/config"
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
func Clone(author, repo string, provider GitRemote) error {
	binPath, err := exec.LookPath("git")
	if err != nil {
		return ErrGitNotInstalled
	}

	cloneURL := provider.CloneURL(author, repo)
	_, err = exec.Command(binPath, "clone", "--filter=blob:none", "--no-checkout", "--depth", "1", cloneURL).Output()
	return err
}

// SparseCheckout configures the current Git repository to only checkout the primary docs directory.
// It returns [ErrGitNotInstalled] if "git" is not located in PATH.
// Note: This function assumes that the repository is already cloned and that the
// current working directory is the root of the repository.
func SparseCheckout() error {
	binPath, err := exec.LookPath("git")
	if err != nil {
		return ErrGitNotInstalled
	}

	_, err = exec.Command(binPath, "sparse-checkout", "init").Output()
	if err != nil {
		return err
	}

	_, err = exec.Command(binPath, "sparse-checkout", "set", config.PrimaryDocsDirName, "README.md").Output()
	if err != nil {
		return err
	}

	_, err = exec.Command(binPath, "checkout").Output()
	return err
}
