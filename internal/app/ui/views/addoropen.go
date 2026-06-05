package views

import (
	"fmt"
	"strings"

	"github.com/inodaf/neoman/internal/app/ui"
	"github.com/inodaf/neoman/internal/app/ui/text"
)

type AddOrOpenViewModel struct {
	Author   string
	Repo     string
	Provider string
	Step     string
	Error    error
	Done     bool
}

// AddOrOpenView renders the progress of adding or opening a project.
//
// Renders:
// │
// ◼  From GitHub - https://github.com/author/repo
// │
// ├  Step 1, 2, ...
// └  Done or Error: <error message>
func AddOrOpenView(model chan AddOrOpenViewModel) {
	ui.Render(model, func(update AddOrOpenViewModel, templ *strings.Builder) {
		templ.WriteString(text.Line("│"))
		templ.WriteString(text.LineBreak())
		templ.WriteString(text.Line("◼  "))

		fmt.Fprint(
			templ,
			text.Colorize("From ", text.Muted, text.Normal),
			text.Colorize(update.Provider, text.Neutral, text.Bold),
			text.Colorize(" - ", text.Neutral, text.Normal),
			text.Colorize(fmt.Sprintf("https://github.com/%s/%s", update.Author, update.Repo), text.Cyan, text.Normal),
			text.LineBreak(),
		)

		templ.WriteString(text.Line("│"))
		templ.WriteString(text.LineBreak())

		if update.Error == nil && !update.Done {
			fmt.Fprint(
				templ,
				text.Line("├  "),
				text.Colorize(update.Step, text.Blue, text.Normal),
				text.LineBreak(),
			)
		} else {
			fmt.Fprint(
				templ,
				text.Line("└  "),
				text.Colorize("Error: ", text.Red, text.Bold),
				text.Colorize(update.Error.Error(), text.Red, text.Italic),
				text.LineBreak(),
			)
		}

		if update.Done {
			fmt.Fprint(
				templ,
				text.Line("└  "),
				text.Colorize("Done - documentation added", text.Purple, text.Bold),
				text.LineBreak(),
			)
		}

		if update.Error != nil {
			return
		}
	})
}
