package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/inodaf/neoman/pkg/config"
)

// ServeIpc serves the IPC Unix Domain Socket (UDS) for communication
// between nman and the nmand (daemon).
func ServeIpc(mux *http.ServeMux) error {
	if err := os.RemoveAll(config.AppSockPath); err != nil {
		slog.Error("unable to remove old socket", "err", err)
		return fmt.Errorf("unable to remove old socket")
	}

	listener, err := net.Listen("unix", config.AppSockPath)
	if err != nil {
		slog.Error("unable to listen to socket", "err", err)
		return fmt.Errorf("unable to listen to socket")
	}

	defer listener.Close()
	slog.Info("Listening to socket", "path", config.AppSockPath)

	if err := http.Serve(listener, mux); err != nil {
		slog.Error("unable to serve socket", "err", err)
		return fmt.Errorf("could not serve Unix Socket")
	}

	return nil
}
