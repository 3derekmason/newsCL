package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

// landing screen state
type homeModel struct {
	choices []string
	cursor  int
}

func newHomeModel() homeModel {
	return homeModel{
		choices: []string{"View Posts", "View Jobs", "Quit"},
	}
}

const (
	homeChoicePosts = iota
	homeChoiceJobs
	homeChoiceQuit
)

// handles input while the home screen is showing
func (m Model) updateHome(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.home.cursor > 0 {
			m.home.cursor--
		}
	case "down", "j":
		if m.home.cursor < len(m.home.choices)-1 {
			m.home.cursor++
		}
	case "enter":
		switch m.home.cursor {
		case homeChoicePosts:
			return m.enterList(sectionPosts)
		case homeChoiceJobs:
			return m.enterList(sectionJobs)
		case homeChoiceQuit:
			return m, tea.Quit
		}
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) viewHome() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Hacker News, in your terminal"))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Browse the current top posts and jobs from Hacker News."))
	b.WriteString("\n\n")

	for i, choice := range m.home.choices {
		b.WriteString(renderChoice(choice, i == m.home.cursor))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("[↑/↓] navigate  [enter] select   [q] quit"))

	return boxStyle.Render(b.String())
}

// renders a single "Label" menu row, highlighted when selected
func renderChoice(label string, selected bool) string {
	if selected {
		return pointerStyle.Render("› ") + selectedItemStyle.Render(label)
	}
	return "  " + itemStyle.Render(label)
}
