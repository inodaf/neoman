package usecase

import (
	"github.com/inodaf/neoman/internal/daemon/repo"
	"github.com/inodaf/neoman/pkg/git"
)

func NewUseCase(docsRepository repo.DocsRepository, gitRemoteClient git.GitRemote) *UseCase {
	return &UseCase{docsRepository: docsRepository, gitRemoteClient: gitRemoteClient}
}

type UseCase struct {
	docsRepository  repo.DocsRepository
	gitRemoteClient git.GitRemote
}
