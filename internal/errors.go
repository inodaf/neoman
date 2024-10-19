package internal

import "errors"

var ErrGetWd = errors.New("neoman: Could not get current working directory")
var ErrReadDocsDir = errors.New("neoman: No 'docs/' or 'manual/' in this directory")
