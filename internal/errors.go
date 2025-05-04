package internal

import "errors"

var (
	ErrGetWd              = errors.New("neoman: Could not get current working directory")
	ErrReadDocsDir        = errors.New("neoman: No 'docs/' directory in this workspace")
	ErrEmptyDocsDir       = errors.New("neoman: Docs directory for '%s' is empty")
	ErrNotAGitRepository  = errors.New("neoman: Directory '%s' is not a Git repository")
	ErrAccessRegistryDir  = errors.New("neoman: Could not access documentation registry")
	ErrAlreadyRegistered  = errors.New("neoman: A project of same name is already registered")
	ErrCreateLocalDocsDir = errors.New("neoman: Could not create 'local' registry directory for project")
	ErrDoubleDashedWdName = errors.New("neoman: Project name starting with '--' as in '%s' is not allowed")
)
