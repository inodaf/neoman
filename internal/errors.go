package internal

import "errors"

var (
	ErrGetWd              = errors.New("could not get current working directory")
	ErrReadDocsDir        = errors.New("no 'docs/' directory in this workspace")
	ErrEmptyDocsDir       = errors.New("docs directory for '%s' is empty")
	ErrNotAGitRepository  = errors.New("directory '%s' is not a Git repository")
	ErrDoubleDashedWdName = errors.New("project name starting with '--' as in '%s' is not allowed")
)
