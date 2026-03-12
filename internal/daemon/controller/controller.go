package controller

import (
	"net/http"

	"github.com/inodaf/neoman/internal/daemon/usecase"
)

func NewHttpController(useCase *usecase.UseCase) *http.ServeMux {
	mux := http.NewServeMux()
	ctrl := &controller{useCase: useCase}
	
	mux.HandleFunc("POST /add/{author}/{repo}", ctrl.AddDocs)
	
	return mux
}

type controller struct {
	useCase *usecase.UseCase
}
