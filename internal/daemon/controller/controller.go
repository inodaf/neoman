package controller

import (
	"net/http"

	"github.com/inodaf/neoman/internal/daemon/usecase"
	"github.com/inodaf/neoman/internal/daemon/worker"
)

func NewHttpController(useCase *usecase.UseCase) *http.ServeMux {
	mux := http.NewServeMux()
	ctrl := &controller{useCase: useCase}

	mux.HandleFunc("POST /add/{author}/{repo}", ctrl.AddDocs)
	mux.HandleFunc("GET /ping", ctrl.Ping)

	return mux
}

type controller struct {
	useCase *usecase.UseCase
	worker  *worker.Worker
}
