package command

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func (c *Command) ListAuthorDocs(ctx context.Context, author string) {
	output, err := c.listAuthorDocs(ctx, author)
	if err != nil {
		fmt.Printf("neoman: %s.\n", err.Error())
		os.Exit(1)
		return
	}

	fmt.Print(output)
	os.Exit(0)
}

func (c *Command) listAuthorDocs(ctx context.Context, author string) (string, error) {
	resource := url.URL{Host: "unix", Scheme: "http", Path: fmt.Sprintf("/docs/%s", author)}

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
		return "", fmt.Errorf("No documentations found for author '%s'", author)
	}

	if res.StatusCode == http.StatusBadRequest {
		return "", ErrListAuthorDocsBadRequest
	}

	if res.StatusCode != http.StatusOK {
		return "", ErrListAuthorDocsUnexpected
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

var (
	ErrListAuthorDocsBadRequest = fmt.Errorf("Author is required")
	ErrListAuthorDocsUnexpected = fmt.Errorf("Unable to list documentations for author")
)
