package driving

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CLIUserPrompter struct{}

func (CLIUserPrompter) ConfirmTrust(authorName string) (bool, error) {
	fmt.Printf("Do you trust owner '%s'? (y/n): ", authorName)

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false, scanner.Err()
	}

	input := strings.ToLower(scanner.Text())
	return input == "y" || input == "yes", nil
}
