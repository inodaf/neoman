package usecases

import (
	"github.com/inodaf/neoman/internal2/ports"
)

type ReadLocalDocs struct {
	ProjectRegistry ports.ProjectRegistry
}

func (rl *ReadLocalDocs) Execute(projectName string, projectPath string) error {
	return rl.ProjectRegistry.AddEntry(ports.RegistryEntry{
		Scope:       ports.LocalScope,
		Project:     projectName,
		ProjectPath: projectPath,
	})
}
