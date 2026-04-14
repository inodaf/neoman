package command

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/inodaf/neoman/pkg/config"
)

func New() *Command {
	daemonClient := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial("unix", config.AppSockPath)
			},
		},
	}

	return &Command{daemonHttpClient: daemonClient}
}

type Command struct {
	daemonHttpClient http.Client
}

func (c *Command) parseAuthorAndRepo(arg string) (author string, repo string, err error) {
	var resourcePattern = regexp.MustCompile(`[^a-zA-Z0-9-_.\s\/]`)
	arg = strings.ToLower(strings.TrimSpace(arg))
	separatorCount := strings.Count(arg, "/")

	if separatorCount > 1 || resourcePattern.MatchString(arg) {
		return "", "", fmt.Errorf("Invalid argument. Must be 'org/repo' or 'author/repo'")
	}

	authorWithRepo := strings.Split(arg, "/")
	author, repo = authorWithRepo[0], authorWithRepo[1]
	err = nil

	return
}
