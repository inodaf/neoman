package internal

import "errors"

var ErrGetWd = errors.New("neoman: Could not get current working directory")
var ErrReadDocsDir = errors.New("neoman: No 'docs/' or 'manual/' in this directory")
var ErrGitNotInstalled = errors.New("neoman: Git is not installed or was not found")
var ErrNotAGitRepository = errors.New("neoman: Directory '%s' is not a Git repository")