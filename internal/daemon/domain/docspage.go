package domain

import (
	"fmt"
	"path"
	"regexp"
	"time"
)

func NewDocsPage(author, repo, title, relPath string) (*DocsPage, error) {
	if author == "" || repo == "" || title == "" {
		return nil, fmt.Errorf("author, repo, and title must be provided")
	}

	page := &DocsPage{
		Title:      title,
		Author:     author,
		Repository: repo,
	}

	err := page.SetRelativePath(relPath)
	if err != nil {
		return nil, err
	}

	return page, nil
}

type DocsPage struct {
	Title          string
	Author         string
	Repository     string
	Content        string
	RelativePath   string
	LastModifiedAt time.Time
	Vector         *map[int]float32
}

func (p *DocsPage) SetRelativePath(relPath string) error {
	var pattern, err = regexp.Compile(`\.md(x)?$`)
	if err != nil {
		return fmt.Errorf("unable to check file extension")
	}

	if pattern.MatchString(path.Ext(relPath)) == false {
		return fmt.Errorf("file '%s' is not a Markdown (.md or .mdx)", relPath)
	}

	p.RelativePath = relPath
	return nil
}

func (p *DocsPage) SetLastModifiedAt(newDate time.Time) error {
	if newDate.After(p.LastModifiedAt) == false || newDate.Equal(p.LastModifiedAt) == false {
		return fmt.Errorf("new date must be after the current last modified date")
	}
	
	p.LastModifiedAt = newDate
	return nil
}