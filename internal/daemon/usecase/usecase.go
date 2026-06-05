package usecase

import (
	"github.com/inodaf/neoman/internal/daemon/repo"
	"github.com/inodaf/neoman/pkg/git"
)

func NewUseCase(
	docsRepository repo.DocsRepository,
	docsPageRepository repo.DocsPageRepository,
	gitHost git.GitRemote,
	sourceRegistry repo.ContentSource,
	jobRepository repo.JobRepository,

) *UseCase {
	return &UseCase{
		docsRepository:     docsRepository,
		docsPageRepository: docsPageRepository,
		sourceRegistry:     sourceRegistry,
		gitHost:            gitHost,
		jobRepository:      jobRepository,
	}
}

type UseCase struct {
	docsRepository     repo.DocsRepository
	docsPageRepository repo.DocsPageRepository
	sourceRegistry     repo.ContentSource
	gitHost            git.GitRemote
	jobRepository      repo.JobRepository
}
