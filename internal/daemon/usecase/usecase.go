package usecase

import (
	"github.com/inodaf/neoman/internal/daemon/repo"
	"github.com/inodaf/neoman/pkg/git"
)

func NewUseCase(
	docsRepository repo.DocsRepository,
	gitRemoteClient git.GitRemote,
	sourceRegistry repo.SourceRegistry,
) *UseCase {
	return &UseCase{
		docsRepository:  docsRepository,
		sourceRegistry:  sourceRegistry,
		gitRemoteClient: gitRemoteClient,
	}
}

type UseCase struct {
	docsRepository  repo.DocsRepository
	sourceRegistry  repo.SourceRegistry
	gitRemoteClient git.GitRemote
}
