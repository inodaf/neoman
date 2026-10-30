package command

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func (c *Command) ListAllDocs(ctx context.Context) {
	output, err := c.listAllDocs(ctx)
	if err != nil {
		fmt.Printf("neoman: %s.\n", err.Error())
		os.Exit(1)
		return
	}

	fmt.Print(output)
	os.Exit(0)
}

func (c *Command) listAllDocs(ctx context.Context) (string, error) {
	resource := url.URL{Host: "unix", Scheme: "http", Path: "/docs"}

	req, err := http.NewRequestWithContext(ctx, "GET", resource.String(), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/yaml")

	res, err := c.daemonClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return "", ErrListAllDocsEmpty
	}

	if res.StatusCode != http.StatusOK {
		return "", ErrListAllDocsUnexpected
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

var (
	ErrListAllDocsEmpty      = fmt.Errorf("no documentations found. add by running 'nman author/repo'")
	ErrListAllDocsUnexpected = fmt.Errorf("unable to list documentations")
)
