package browser

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/inodaf/neoman/packages/config"
)

func Open(proj string) error {
	url := fmt.Sprintf("http://%s/%s", config.AppHostName, proj)
	browser := os.Getenv("BROWSER")

	if len(browser) == 0 {
		fmt.Printf("Open %s in your browser.\n", url)
		return nil
	}

	fmt.Printf("Opening %s in your browser.\n", url)
	cmd := exec.Command(browser, url)

	err := cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()
}
