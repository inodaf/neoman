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

	proj := path.Base(wd)

	if !internal.IsGitRepository() {
		fmt.Printf(internal.ErrNotAGitRepository.Error(), proj)
		return
	}

	if internal.IsAlreadyRegistered(proj) {
		 // TODO: Handle text-based rendering (terminal only)
		fmt.Printf("Opening https://neoman.local/%s in your browser.\n", proj)
		return
	}

	// TODO: Change to os.Stat() ?
	// TODO: Add safe checks
	_, err = os.ReadDir(path.Join(wd, internal.PrimaryDocsDirName))
	if err != nil {
		_, err = os.ReadDir(path.Join(wd, internal.AlternateDocsDirName))
		if err != nil {
			fmt.Println(internal.ErrReadDocsDir)
			return
		}
	}

	err = internal.AddSymlinkToRegistry(proj, wd)
	if err != nil {
		fmt.Printf("neoman: Could not link working directory docs on registry")
		return
	}
}

// TODO: Link WD into "Registry"
