package internal

import (
	"errors"
	"log"
	"os/exec"
)

var ErrGitNotInstalled = errors.New("neoman: Git is not installed or was not found")

// IsGitRepository checks if the current working directory
// is a valid Git repository by checking the "git status" command
// output. It Exit(1) if the "git" binary was not located in PATH.
func IsGitRepository() bool {
	git, err := exec.LookPath("git")
	if err != nil {
		log.Fatalln(ErrGitNotInstalled.Error())
	}

	_, err = exec.Command(git, "status").Output()
	return err == nil
}
