package daemon

import (
	"net"

	"github.com/inodaf/neoman/internal"
)

func Ping() error {
	addr, err := net.ResolveUnixAddr("unix", internal.AppSockAddr)
	if err != nil {
		return err
	}

	conn, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("ping")); err != nil {
		return err
	}

	return nil
}