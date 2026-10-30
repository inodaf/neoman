package command

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
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

	return &Command{daemonClient: daemonClient}
}

type Command struct {
	daemonClient http.Client
}

func (c *Command) parseAuthorAndRepo(arg string) (author string, repo string, err error) {
	resourcePattern := regexp.MustCompile(`[^a-zA-Z0-9-_.\s\/]`)
	arg = strings.ToLower(strings.TrimSpace(arg))
	separatorCount := strings.Count(arg, "/")

	if separatorCount > 1 || resourcePattern.MatchString(arg) {
		return "", "", fmt.Errorf("invalid argument. Must be 'org/repo' or 'author/repo'")
	}

	authorWithRepo := strings.Split(arg, "/")
	author, repo = authorWithRepo[0], authorWithRepo[1]
	err = nil

	return
}

// SocketClientPing tries to ping the daemon through the Unix socket.
// If the daemon is not running, it returns an error.
func (c *Command) SocketClientPing() error {
	resource := url.URL{Host: "unix", Scheme: "http", Path: "/ping"}
	_, err := c.daemonClient.Get(resource.String())

	return err
}
