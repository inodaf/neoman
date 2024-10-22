package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/inodaf/neoman/internal"
)

func OpenFromWD() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(internal.ErrGetWd)
		return
	}

	proj := path.Base(wd)

	if strings.HasPrefix(proj, "--") {
		fmt.Printf(internal.ErrDoubleDashedWdName.Error(), proj)
		return
	}

	if !internal.IsGitRepository() {
		fmt.Printf(internal.ErrNotAGitRepository.Error(), proj)
		return
	}

	if internal.IsAlreadyRegistered(proj) {
		err = viewDocs(proj)
		if err != nil {
			fmt.Println("neoman: Could not display the docs for this project")
		}
		return
	}

	// TODO: Add safe checks
	docsDir, err := os.ReadDir(path.Join(wd, internal.PrimaryDocsDirName))
	if err != nil {
		fmt.Println(internal.ErrReadDocsDir)
		return
	}

	if len(docsDir) == 0 {
		fmt.Printf(internal.ErrEmptyDocsDir.Error(), proj)
		return
	}

	err = internal.AddSymlinkToRegistry(proj, wd)
	if err != nil {
		fmt.Printf("neoman: Could not link working directory docs on registry")
		return
	}

	viewDocs(proj)
}

func viewDocs(proj string) error {
	// TODO: Handle text-based rendering (terminal only)
	defaultBrowser := os.Getenv("BROWSER")
	if len(defaultBrowser) != 0 {
		fmt.Printf("Opening https://neoman.local/%s in your browser.\n", proj)
		return exec.Command(defaultBrowser, fmt.Sprintf("https://neoman.local/%s", proj)).Start()
	} else {
		fmt.Printf("Open https://neoman.local/%s in your browser.\n", proj)
	}

	return nil
}
