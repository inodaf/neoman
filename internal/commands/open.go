package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

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

	proj := path.Base(wd)

	if strings.HasPrefix(proj, "--") {
		fmt.Printf(internal.ErrDoubleDashedWdName.Error(), proj)
		return
	}

	ok, err := git.IsRepository()
	if err != nil {
		fmt.Println(err)
		return
	}
	if !ok {
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

	if err = internal.AddLocalEntryToRegistry(proj, wd); err != nil {
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

	ownerWithRepo := strings.Split(proj, "/")
	if ok := daemon.IPC.IsAccountTrusted(ownerWithRepo[0]); !ok {
		if confirmTrust(ownerWithRepo[0]) {
			daemon.IPC.TrustAccount(ownerWithRepo[0])
		} else {
			fmt.Printf("neoman:	'%s' is not trusted. Stopping.\n", ownerWithRepo[0])
			return
		}
	}

	// TODO: Handle text-based rendering (terminal only) based on user preferences.
	if internal.IsAlreadyRegistered(proj) || internal.IsAlreadyRegistered(proj, "remote") {
		if err := openInBrowser(proj); err != nil {
			fmt.Println("neoman: Could not display the docs for this project")
		}
		return
	}

	fmt.Printf("neoman:	Repository '%s' is not registered. Trying Git remotes...\n", proj)
	err := internal.FetchDocs(ownerWithRepo[0], ownerWithRepo[1])
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
	url := fmt.Sprintf("https://%s/%s", internal.AppHostName, proj)

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
