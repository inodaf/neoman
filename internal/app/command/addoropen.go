package command

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func (c *Command) AddOrOpen(ctx context.Context, arg string) {
	author, repo, err := c.parseAuthorAndRepo(arg)
	if err != nil {
		fmt.Printf("neoman: %s.\n", err.Error())
		return
	}

	err = c.addDocs(ctx, author, repo)
	if err != nil && errors.Is(err, ErrConflict) {
		fmt.Printf("neoman: %s.\n", err.Error())
		os.Exit(1)
		return
	} else if err != nil {
		fmt.Printf("neoman: %s.\n", err.Error())
		os.Exit(1)
		return
	}

	os.Exit(0)
}

func (c *Command) addDocs(ctx context.Context, author, repo string) error {
	resource := url.URL{Host: "unix", Scheme: "http", Path: fmt.Sprintf("/add/%s/%s", author, repo)}

	req, err := http.NewRequestWithContext(ctx, "POST", resource.String(), nil)
	if err != nil {
		return err
	}

	res, err := c.daemonHttpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusConflict {
		return ErrConflict
	}

	if res.StatusCode == http.StatusNoContent {
		return ErrNoDocsDir
	}

	if res.StatusCode != http.StatusAccepted {
		return ErrUnexpected
	}

	return nil
}

var (
	ErrConflict   = fmt.Errorf("Docs already added")
	ErrNoDocsDir  = fmt.Errorf("Repo does not have a 'docs/' directory")
	ErrUnexpected = fmt.Errorf("Unable to add docs")
)
