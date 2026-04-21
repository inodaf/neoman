package command

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func (c *Command) ViewPage(ctx context.Context, project string, docPath string) {
	author, repo, err := c.parseAuthorAndRepo(project)
	if err != nil {
		fmt.Printf("neoman: %s.\n", err.Error())
		os.Exit(1)
		return
	}

	output, err := c.viewDocPage(ctx, author, repo, docPath)
	if err != nil {
		fmt.Printf("neoman: %s.\n", err.Error())
		os.Exit(1)
		return
	}

	fmt.Print(output)
	os.Exit(0)
}

func (c *Command) viewDocPage(ctx context.Context, author, repo, docPath string) (string, error) {
	// URL-encode the docPath to handle spaces and special characters
	encodedPath := url.PathEscape(docPath)

	resource := url.URL{
		Host:   "unix",
		Scheme: "http",
		Path:   fmt.Sprintf("/docs/%s/%s/pages/%s", author, repo, encodedPath),
	}

	req, err := http.NewRequestWithContext(ctx, "GET", resource.String(), nil)
	if err != nil {
		return "", err
	}

	res, err := c.daemonHttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return "", ErrViewPageNotFound
	}

	if res.StatusCode != http.StatusOK {
		return "", ErrViewPageUnexpected
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

var (
	ErrViewPageNotFound   = fmt.Errorf("Page not found. Check the path and try again")
	ErrViewPageUnexpected = fmt.Errorf("Unable to view page")
)
