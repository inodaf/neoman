package main

import (
	"sync"

	"github.com/inodaf/neoman/internal/management"
)

func main() {
	db, err := management.NewSQLiteDatabase()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		management.SocketServeIPC(db)
	}()

	go func() {
		defer wg.Done()
		management.SocketServeTCP(db)
	}()

	wg.Wait()
}
