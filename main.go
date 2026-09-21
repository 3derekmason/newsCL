// Command newscl is a terminal UI for browsing current Hacker News top
// posts and jobs, built with Bubble Tea UI framework.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"newscl/internal/ui"
)

func main() {
	p := tea.NewProgram(ui.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
