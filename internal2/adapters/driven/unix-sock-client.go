package driven

import (
  "context"
  "github.com/inodaf/neoman/packages/config"
  "net"
  "net/http"
  "net/url"
)

var UnixSockClient = http.Client{
  Transport: &http.Transport{
    DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
      return net.Dial("unix", config.AppSockAddr)
    },
  },
}

func PingUnixSock() error {
	resource := url.URL{Host: "unix", Scheme: "http", Path: "/ping"}
	_, err := UnixSockClient.Get(resource.String())

	return err
}