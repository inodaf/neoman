package worker

import (
	"github.com/inodaf/neoman/internal/daemon/repo"
)

func NewWorker(
	docsPageRepository repo.DocsPageRepository,
	sourceRegistry repo.SourceRegistry,
	docsRepository repo.DocsRepository,
) *Worker {
	return &Worker{
		docsPageRepository: docsPageRepository,
		sourceRegistry:     sourceRegistry,
		docsRepository:     docsRepository,
	}
}

type Worker struct {
	sourceRegistry     repo.SourceRegistry
	docsPageRepository repo.DocsPageRepository
	docsRepository     repo.DocsRepository
}
