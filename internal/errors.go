package internal

import "errors"

var (
	ErrGetWd              = errors.New("neoman: Could not get current working directory")
	ErrReadDocsDir        = errors.New("neoman: No 'docs/' or 'manual/' in this directory")
	ErrGitNotInstalled    = errors.New("neoman: Git is not installed or was not found")
	ErrNotAGitRepository  = errors.New("neoman: Directory '%s' is not a Git repository")
	ErrAccessRegistryDir  = errors.New("neoman: Could not access documentation registry")
	ErrAlreadyRegistered  = errors.New("neoman: A project of same name is already registered")
	ErrCreateLocalDocsDir = errors.New("neoman: Could not create 'local' registry directory for project")
)
