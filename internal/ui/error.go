package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

// the state for the error screen
type errModel struct {
	message string  // the error message
	section section // the section to retry
}

// switches to the error screen with the given message, ready to
// retry the given section on request.
func (m Model) enterError(err error, sec section) Model {
	m.screen = screenError
	m.err = errModel{message: err.Error(), section: sec}
	return m
}

func (m Model) updateError(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "r":
		return m.enterList(m.err.section)
	case "left", "q":
		m.screen = screenHome
	case "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) viewError() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Oops! Something went wrong"))
	b.WriteString("\n\n")
	b.WriteString(errorStyle.Render(m.err.message))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("[r] retry   [←] back to home   [q] quit"))

	return boxStyle.Width(60).Render(b.String())
}
