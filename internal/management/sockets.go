package management

import (
	"context"
	"net"
	"net/http"

	"github.com/inodaf/neoman/pkg/config"
)

var UnixSockClient http.Client = http.Client{
	Transport: &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("unix", config.AppSockPath)
		},
	},
}
