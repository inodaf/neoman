package main

import (
	"github.com/inodaf/neoman/internal"
	"github.com/inodaf/neoman/internal/daemon"
)

func main() {
	daemon.ServeIPC(internal.AppSockAddr)
}
