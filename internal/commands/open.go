package commands

import (
	"fmt"
	"github.com/inodaf/neoman/internal"
	"os"
	"path"
)

func OpenFromWD() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(internal.ErrGetWd)
		return
	}

	projectName := path.Base(wd)

	if !internal.IsGitRepository() {
		fmt.Printf(internal.ErrNotAGitRepository.Error(), projectName)
		return
	}

	docsDir, err := os.ReadDir(path.Join(wd, internal.PrimaryDocsDirName))
	if err != nil {
		docsDir, err = os.ReadDir(path.Join(wd, internal.AlternateDocsDirName))
		if err != nil {
			fmt.Println(internal.ErrReadDocsDir)
			return
		}
	}

	fmt.Printf("Manual files from project '%s':\n\n", projectName)
	for _, v := range docsDir {
		fmt.Println(v.Name())
	}
}

// @TODO: Link WD into "Registry"
// @TODO: Skip WD Link if already in "Registry"