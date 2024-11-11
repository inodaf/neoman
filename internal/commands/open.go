package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
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
		if err = openInBrowser(proj); err != nil {
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

	if err = internal.AddSymlinkToRegistry(proj, wd); err != nil {
		fmt.Printf("neoman: Could not link working directory docs on registry")
		return
	}

	openInBrowser(proj)
}

func OpenFromName(proj string) {
	re := regexp.MustCompile(`[^a-z0-9-_.\s\/]`)

	if strings.Count(proj, "/") > 1 || re.Match([]byte(proj)) {
		fmt.Println("neoman: Invalid argument. Must be 'repo' or 'org/repo'")
		return
	}

	if internal.IsAlreadyRegistered(proj) || internal.IsAlreadyRegistered(proj, "remote") {
		// TODO: Handle text-based rendering (terminal only)
		// based on user preferences.
		if err := openInBrowser(proj); err != nil {
			fmt.Println("neoman: Could not display the docs for this project")
		}
		return
	}

	fmt.Printf("neoman: Project '%s' not registered, trying Git remotes\n", proj)
}

func openInBrowser(proj string) error {
	browser := os.Getenv("BROWSER")
	if len(browser) != 0 {
		fmt.Printf("Opening https://neoman.local/%s in your browser.\n", proj)
		cmd := exec.Command(browser, fmt.Sprintf("https://neoman.local/%s", proj))
		if err := cmd.Start(); err != nil {
			return err
		}
		return cmd.Wait()
	}

	fmt.Printf("Open https://neoman.local/%s in your browser.\n", proj)
	return nil
}
