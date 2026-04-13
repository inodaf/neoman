package controller

import (
	"net/http"

	"github.com/inodaf/neoman/internal/daemon/repo"
	"github.com/inodaf/neoman/internal/daemon/usecase"
)

func NewHttpController(useCase *usecase.UseCase, jobRepository repo.JobRepository) *http.ServeMux {
	mux := http.NewServeMux()
	ctrl := &controller{useCase: useCase, jobRepository: jobRepository}

	mux.HandleFunc("POST /add/{author}/{repo}", ctrl.AddDocs)
	mux.HandleFunc("GET /ping", ctrl.Ping)
	mux.HandleFunc("GET /internal/jobs", ctrl.GetJobsStatus)
	mux.HandleFunc("GET /docs/{author}/{repo}/pages", ctrl.ListDocPages)

	return mux
}

type controller struct {
	useCase       *usecase.UseCase
	jobRepository repo.JobRepository
}
