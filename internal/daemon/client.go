package daemon

import (
	"context"
	"net"
	"net/http"
	"net/url"

	"github.com/inodaf/neoman/internal"
)

type unixSockIPC struct {
	client http.Client
}

func (ipc *unixSockIPC) Ping() error {
	resource := url.URL{Host: "unix", Scheme: "http", Path: "/ping"}
	if _, err := ipc.client.Get(resource.String()); err != nil {
		return err
	}
	return nil
}

func newUnixSockClient(sockAddr string) *unixSockIPC {
	return &unixSockIPC{client: http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return net.Dial("unix", sockAddr)
			},
		},
	}}
}

var IPC = newUnixSockClient(internal.AppSockAddr)
