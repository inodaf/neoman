package text

import "fmt"

type Color string

const (
	Black   Color = "30m"
	Red     Color = "31m"
	Green   Color = "32m"
	Yellow  Color = "33m"
	Blue    Color = "34m"
	Purple  Color = "35m"
	Cyan    Color = "36m"
	Neutral Color = "37m"
	Muted   Color = "2m"
	Reset   Color = "0m"
)

type Style string

const (
	Normal    Style = "0"
	Bold      Style = "1"
	Italic    Style = "3"
	Underline Style = "4"
)

func Colorize(text string, color Color, style Style) string {
	return fmt.Sprintf("\033[%s;%s%s\033[%s", style, color, text, Reset)
}
