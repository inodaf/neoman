package main

import (
	"context"
	"sync"

	"github.com/inodaf/neoman/internal/daemon/controller"
	"github.com/inodaf/neoman/internal/daemon/repo"
	"github.com/inodaf/neoman/internal/daemon/scheduler"
	"github.com/inodaf/neoman/internal/daemon/usecase"
	"github.com/inodaf/neoman/internal/daemon/worker"
	"github.com/inodaf/neoman/internal/management"
	"github.com/inodaf/neoman/pkg/git"
)

func main() {
	db, err := InitDB()
	if err != nil {
		panic(err)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ghClient := git.NewGitHubClient()
	jobRepository := repo.NewJobRepository(db)
	docsRepository := repo.NewDocsRepository(db)
	docsPageRepository := repo.NewDocsPageRepository(db)
	fsSourceRegistry := repo.NewFsSourceRegistry(ghClient)

	worker := worker.NewWorker(docsPageRepository, fsSourceRegistry, docsRepository, jobRepository)
	useCase := usecase.NewUseCase(docsRepository, docsPageRepository, ghClient, fsSourceRegistry, jobRepository)
	scheduler := scheduler.NewScheduler(worker)

	mux := controller.NewHttpController(useCase, jobRepository)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		if err := scheduler.Run(ctx); err != nil && err != context.Canceled {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		err := ServeIpc(mux)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		management.SocketServeTCP(db)
	}()

	wg.Wait()
}
