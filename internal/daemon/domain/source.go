package domain

type ContentSource string

const (
	SourceGitHub ContentSource = "github"
	SourceGitLab ContentSource = "gitlab"
)
