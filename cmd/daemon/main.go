package main

import (
	"github.com/inodaf/neoman/internal"
	"github.com/inodaf/neoman/internal/daemon"
)

func main() {
	db, err := internal.NewSQLiteDatabase()
	if err != nil {
		panic(err)
	}

	internal.DB = db
	daemon.ServeIPC(internal.AppSockAddr)
}
