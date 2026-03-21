package main

import (
	"sync"

	"github.com/inodaf/neoman/internal/daemon/controller"
	"github.com/inodaf/neoman/internal/daemon/repo"
	"github.com/inodaf/neoman/internal/daemon/usecase"
	"github.com/inodaf/neoman/internal/management"
	"github.com/inodaf/neoman/pkg/git"
)

func main() {
	dbConn, err := InitDB()
	if err != nil {
		panic(err)
	}

	ghClient := git.NewGitHubClient()
	fsSourceRegistry := repo.NewFsSourceRegistry(ghClient)
	useCase := usecase.NewUseCase(nil, ghClient, fsSourceRegistry)
	mux := controller.NewHttpController(useCase)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := ServeIpc(mux)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		management.SocketServeTCP(dbConn)
	}()

	wg.Wait()
}
