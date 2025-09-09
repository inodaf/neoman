package browser

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/inodaf/neoman/packages/config"
)

var fallbackMessage = "Open %s in your browser.\n"

func Open(proj string) {
	url := fmt.Sprintf("http://%s/%s", config.AppHostName, proj)
	browser := os.Getenv("BROWSER")

	if len(browser) == 0 {
		fmt.Printf(fallbackMessage, url)
		return
	}

	cmd := exec.Command(browser, url)
	err := cmd.Start()

	if err != nil {
		fmt.Printf(fallbackMessage, url)
		return
	}

	cmd.Wait()
}
