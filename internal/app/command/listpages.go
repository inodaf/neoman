package command

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func (c *Command) ListPages(ctx context.Context, arg string) {
	author, repo, err := c.parseAuthorAndRepo(arg)
	if err != nil {
		fmt.Printf("neoman: %s.\n", err.Error())
		os.Exit(1)
		return
	}

	output, err := c.listDocPages(ctx, author, repo)
	if err != nil {
		fmt.Printf("neoman: %s.\n", err.Error())
		os.Exit(1)
		return
	}

	fmt.Print(output)
	os.Exit(0)
}

func (c *Command) listDocPages(ctx context.Context, author, repo string) (string, error) {
	resource := url.URL{Host: "unix", Scheme: "http", Path: fmt.Sprintf("/docs/%s/%s/pages", author, repo)}

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
		return "", ErrListDocsNotFound
	}

	if res.StatusCode == http.StatusNoContent {
		return "", ErrListNoPagesFound
	}

	if res.StatusCode != http.StatusOK {
		return "", ErrListUnexpected
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

var (
	ErrListDocsNotFound = fmt.Errorf("Documentation not found. Add by running 'neoman author/repo'")
	ErrListNoPagesFound = fmt.Errorf("Documentation is empty")
	ErrListUnexpected   = fmt.Errorf("Unable to show documentation pages")
)
