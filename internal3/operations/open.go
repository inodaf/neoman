package operations

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/inodaf/neoman/internal3/management"
	"github.com/inodaf/neoman/packages/browser"
	"github.com/inodaf/neoman/packages/config"
	"github.com/inodaf/neoman/packages/git"
)

func OpenFromCurrentDirectory() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(management.ErrGetWd)
		return
	}

	project := path.Base(wd)

	if strings.HasPrefix(project, "--") {
		fmt.Printf(management.ErrDoubleDashedWdName.Error(), project)
		return
	}

	if ok, err := git.IsRepository(); !ok || err != nil {
		fmt.Printf(management.ErrNotAGitRepository.Error(), project)
		return
	}

	newProjectEntry := management.RegistryEntry{
		Scope:       management.RegistryTypeLocal,
		Project:     project,
		ProjectPath: wd,
	}

	if management.RegistryHasEntry(newProjectEntry) {
		if err = browser.Open(project); err != nil {
			fmt.Println("neoman: Could not display the docs for this project")
		}
		return
	}

	// TODO: Add safe checks
	docsDir, err := os.ReadDir(path.Join(wd, config.PrimaryDocsDirName))
	if err != nil {
		fmt.Println(management.ErrReadDocsDir)
		return
	}

	if len(docsDir) == 0 {
		fmt.Printf(management.ErrEmptyDocsDir.Error(), project)
		return
	}

	if err = management.RegistryAddEntry(newProjectEntry); err != nil {
		fmt.Printf("neoman: Could not link working directory docs on registry")
		return
	}

	browser.Open(project)
}

func OpenFromRepositoryName(proj string) {
	re := regexp.MustCompile(`[^a-zA-Z0-9-_.\s\/]`)
	if strings.Count(proj, "/") > 1 || re.Match([]byte(proj)) {
		fmt.Println("neoman: Invalid argument. Must be 'repo' or 'org/repo'")
		return
	}

	proj = strings.ToLower(proj)
	authorWithRepo := strings.Split(proj, "/")
	author, repo := authorWithRepo[0], authorWithRepo[1]

	// TODO: Handle text-based rendering (terminal only) based on user preferences.
	if management.RegistryHasEntry(management.RegistryEntry{
		Scope:   management.RegistryTypeLocal,
		Project: proj,
	}) || management.RegistryHasEntry(management.RegistryEntry{
		Scope:   management.RegistryTypeRemote,
		Project: proj,
		Owner:   author,
	}) {
		if err := browser.Open(proj); err != nil {
			fmt.Println("neoman: Could not display the docs for this project")
		}
		return
	}

	err := FetchDocs(author, repo)
	if err != nil {
		fmt.Print(err)
		return
	}

	if err := browser.Open(proj); err != nil {
		fmt.Println("neoman: Could not display the docs for this project")
	}
}
