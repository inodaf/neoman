package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/inodaf/neoman/internal"
	"github.com/inodaf/neoman/internal/daemon"
	"github.com/inodaf/neoman/internal/git"
)

func OpenFromWD() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(internal.ErrGetWd)
		return
	}

	project := path.Base(wd)

	if strings.HasPrefix(project, "--") {
		fmt.Printf(internal.ErrDoubleDashedWdName.Error(), project)
		return
	}

	if ok, err := git.IsRepository(); !ok ||  err != nil {
		fmt.Printf(internal.ErrNotAGitRepository.Error(), project)
		return
	}

	if internal.IsAlreadyRegistered(project) {
		if err = openInBrowser(project); err != nil {
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
		fmt.Printf(internal.ErrEmptyDocsDir.Error(), project)
		return
	}

	if err = internal.AddLocalEntryToRegistry(project, wd); err != nil {
		fmt.Printf("neoman: Could not link working directory docs on registry")
		return
	}

	openInBrowser(project)
}

func OpenFromName(proj string) {
	re := regexp.MustCompile(`[^a-zA-Z0-9-_.\s\/]`)
	if strings.Count(proj, "/") > 1 || re.Match([]byte(proj)) {
		fmt.Println("neoman: Invalid argument. Must be 'repo' or 'org/repo'")
		return
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	proj = strings.ToLower(proj)
	ownerWithRepo := strings.Split(proj, "/")
	owner, repo := ownerWithRepo[0], ownerWithRepo[1]

	if ok := daemon.IPC.IsAccountTrusted(owner); len(ownerWithRepo) == 2 && !ok {
		if !confirmTrust(owner) {
			fmt.Printf("neoman:	'%s' is not trusted. Stopping.\n", owner)
			return
		}
		wg.Add(1)
		go func() {
			daemon.IPC.TrustAccount(owner)
			wg.Done()
		}()
	}

	// TODO: Handle text-based rendering (terminal only) based on user preferences.
	if internal.IsAlreadyRegistered(proj) || internal.IsAlreadyRegistered(proj, "remote") {
		if err := openInBrowser(proj); err != nil {
			fmt.Println("neoman: Could not display the docs for this project")
		}
		return
	}

	err := internal.FetchDocs(owner, repo)
	if err != nil {
		fmt.Print(err)
		return
	}

	if err := openInBrowser(proj); err != nil {
		fmt.Println("neoman:	Could not display the docs for this project")
	}
}

func openInBrowser(proj string) error {
	browser := os.Getenv("BROWSER")
	url := fmt.Sprintf("http://%s/%s", internal.AppHostName, proj)

	if len(browser) != 0 {
		fmt.Printf("Opening %s in your browser.\n", url)
		cmd := exec.Command(browser, url)
		if err := cmd.Start(); err != nil {
			return err
		}
		return cmd.Wait()
	}

	fmt.Printf("Open %s in your browser.\n", url)
	return nil
}

func confirmTrust(owner string) bool {
	fmt.Printf("Do you trust owner '%s'? (y/n): ", owner)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	input := strings.ToLower(scanner.Text())
	return input == "y" || input == "yes"
}
