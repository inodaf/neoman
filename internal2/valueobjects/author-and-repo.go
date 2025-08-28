package valueobjects

import (
	"errors"
	"regexp"
	"strings"
)

type AuthorAndRepo string

func (ar AuthorAndRepo) Validate() error {
	re := regexp.MustCompile(`[^a-zA-Z0-9-_.\s/]`)

	if strings.Count(string(ar), "/") > 1 || re.Match([]byte(ar)) {
		return errors.New("invalid argument. Must be 'repo' or 'org/repo'")
	}

	return nil
}

func (ar AuthorAndRepo) Value() (authorName, repo string) {
	authorWithRepo := strings.Split(strings.ToLower(string(ar)), "/")
	return authorWithRepo[0], authorWithRepo[1]
}
