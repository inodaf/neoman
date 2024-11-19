package daemon

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/inodaf/neoman/internal"
)

type unixSockIPC struct {
	client http.Client
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

func (ipc *unixSockIPC) Ping() error {
	resource := url.URL{Host: "unix", Scheme: "http", Path: "/ping"}
	_, err := ipc.client.Get(resource.String())

	return err
}

func (ipc *unixSockIPC) IsAccountTrusted(account string) bool {
	query := url.Values{}
	query.Add("account", account)

	resource := url.URL{
		Scheme:   "http",
		Host:     "unix",
		Path:     "/trust",
		RawQuery: query.Encode(),
	}

	resp, err := ipc.client.Get(resource.String())
	if err != nil {
		log.Fatalln("neoman: Could not check account from daemon")
	}

	return resp.StatusCode == http.StatusOK
}

func (ipc *unixSockIPC) TrustAccount(account string) bool {
	query := url.Values{}
	query.Add("account", account)

	resource := url.URL{
		Scheme:   "http",
		Host:     "unix",
		Path:     "/trust",
		RawQuery: query.Encode(),
	}

	resp, err := ipc.client.Post(resource.String(), "text/plain", strings.NewReader(""))
	if err != nil {
		log.Fatalln("neoman: Could not check account from daemon")
	}

	return resp.StatusCode == http.StatusCreated
}
