package operations

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/inodaf/neoman/internal/management"
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
		browser.Open(project)
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

func OpenFromName(proj string) {
	re := regexp.MustCompile(`[^a-zA-Z0-9-_.\s\/]`)
	proj = strings.ToLower(strings.TrimSpace(proj))
	separatorCount := strings.Count(proj, "/")

	if separatorCount > 1 || re.Match([]byte(proj)) {
		fmt.Println("neoman: Invalid argument. Must be 'repo' or 'org/repo'")
		return
	}

	hasLocalEntry := management.RegistryHasEntry(management.RegistryEntry{
		Scope:   management.RegistryTypeLocal,
		Project: proj,
	})

	if separatorCount == 0 && hasLocalEntry {
		browser.Open(proj)
		return
	} else if separatorCount == 0 && !hasLocalEntry {
		return
	}

	authorWithRepo := strings.Split(proj, "/")
	author, repo := authorWithRepo[0], authorWithRepo[1]
	hasRemoteEntry := management.RegistryHasEntry(management.RegistryEntry{
		Scope:   management.RegistryTypeRemote,
		Project: repo,
		Owner:   author,
	})

	// TODO: Handle text-based rendering (terminal only) based on user preferences.
	if hasRemoteEntry {
		browser.Open(proj)
		return
	}

	err := FetchDocs(author, repo)
	if err != nil {
		fmt.Print(err)
		return
	}

	browser.Open(proj)
}
