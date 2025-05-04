package main

import (
	"sync"

	"github.com/inodaf/neoman/internal"
	"github.com/inodaf/neoman/internal/daemon"
)

func main() {
	var wg sync.WaitGroup
	db, err := internal.NewSQLiteDatabase()
	if err != nil {
		panic(err)
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		daemon.ServeSocket(internal.AppSockAddr, db)
	}()

	go func() {
		defer wg.Done()
		daemon.ServeWeb(internal.AppWebAppPort, db)
	}()

	wg.Wait()
}
