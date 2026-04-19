package ui

import (
	"fmt"
	"strings"
)

// Render continuously listens for updates of a Model and renders/re-renders
// a text view in the standard output. For subsequent renders, it clear the previous output
// and recompute the view using [onUpdate].
func Render[T any](model chan T, onUpdate func(update T, templ *strings.Builder)) {
	var templ strings.Builder
	firstRender := true
	lineCount := 0

	for update := range model {
		templ.Reset()

		if firstRender == false {
			fmt.Printf("\033[%dA", lineCount)
		}

		onUpdate(update, &templ)
		output := templ.String()

		firstRender = false
		lineCount = strings.Count(output, "\n")

		fmt.Print(output)
	}
}
