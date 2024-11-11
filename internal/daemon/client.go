package daemon

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/inodaf/neoman/internal"
)

type unixSockIPC struct {
	client http.Client
}

func (ipc *unixSockIPC) Ping() {
	resource := url.URL{Host: "unix", Scheme: "http", Path: "/ping"}

	resp, err := ipc.client.Get(resource.String())
	if err != nil {
		log.Fatal("neoman: could not connect to daemon.")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("neoman: could not read from daemon when pinging.")
	}

	fmt.Println(string(body))
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
