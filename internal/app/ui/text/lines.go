package text

import "fmt"

// Line returns a string that, when printed, will clear the current line and print the provided content.
func Line(content string) string {
	return fmt.Sprintf("\033[2K\r%s", content)
}

func LineBreak() string {
	return "\n"
}