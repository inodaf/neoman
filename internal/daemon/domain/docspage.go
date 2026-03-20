package domain

import (
	"fmt"
	"path"
	"regexp"
	"strings"
	"time"
)

func NewDocsPage(author, repo, content, relPath string) (*DocsPage, error) {
	if author == "" || repo == "" || content == "" {
		return nil, fmt.Errorf("author, repo, and content must be provided")
	}

	page := &DocsPage{
		Author:     author,
		Repository: repo,
		Content:    content,
	}

	err := page.SetRelativePath(relPath)
	if err != nil {
		return nil, err
	}
	
	err = page.SetTitle()
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

func (p *DocsPage) SetTitle() error {
	if p.Content == "" {
		return fmt.Errorf("unable to infer title with no content")
	}
	
	lines := strings.Split(p.Content, "\n")
	for _, line := range lines {
		content, found := strings.CutPrefix(line, "# ")
		if found {
			p.Title = content
			break
		}
	}
	
	if p.RelativePath == "" {
		return fmt.Errorf("unable to infer title with no relative path")
	}
	
	p.Title = strings.TrimSuffix(path.Base(p.RelativePath), path.Ext(p.RelativePath))
	return nil
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