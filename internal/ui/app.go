// ui implements the Bubble Tea application
package ui

import (
	tea "charm.land/bubbletea/v2"

	"newscl/internal/hn"
)

// identifies which "page" of the app is currently showing
type screen int

const (
	screenHome screen = iota
	screenList
	screenModal
	screenError
)

// the root Bubble Tea model
type Model struct {
	client *hn.Client

	width  int
	height int

	screen screen

	home  homeModel
	list  listModel
	modal modalModel
	err   errModel
}

// returns the initial application model, showing the home screen.
func NewModel() Model {
	return Model{
		client: hn.NewClient(),
		screen: screenHome,
		home:   newHomeModel(),
	}
}

// no I/O to kick off until the user picks something from the home menu.
func (m Model) Init() tea.Cmd {
	return nil
}

// dispatches to the handler for whichever screen is currently active.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width, m.height = sizeMsg.Width, sizeMsg.Height
		return m, nil
	}

	switch m.screen {
	case screenHome:
		if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
			return m.updateHome(keyMsg)
		}
	case screenList:
		return m.updateList(msg)
	case screenModal:
		return m.updateModal(msg)
	case screenError:
		if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
			return m.updateError(keyMsg)
		}
	}

	return m, nil
}

// returns the UI for the current screen.
func (m Model) View() tea.View {
	var content string
	switch m.screen {
	case screenHome:
		content = m.viewHome()
	case screenList:
		content = m.viewList()
	case screenModal:
		content = m.viewModal()
	case screenError:
		content = m.viewError()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}
