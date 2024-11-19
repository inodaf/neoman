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

	accountWithRepo := strings.Split(proj, "/")
	if ok := daemon.IPC.IsAccountTrusted(accountWithRepo[0]); !ok {
		if askToTrust(accountWithRepo[0]) {
			daemon.IPC.TrustAccount(accountWithRepo[0])
		} else {
			fmt.Printf("neoman: User '%s' not trusted.\n", accountWithRepo[0])
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

	fmt.Printf("neoman: Project '%s' not registered, trying Git remotes\n", proj)
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

func askToTrust(account string) bool {
	fmt.Printf("Remote account '%s' is not trusted. Do you want to trust it? (y/n): ", account)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	input := strings.ToLower(scanner.Text())
	return input == "y" || input == "yes"
}
