package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/pkg/browser"
)

// opens url in the user's default browser
func openURL(url string) error {
	return browser.OpenURL(url)
}

// reports the outcome of an openURLCmd back to Update.
type browserResultMsg struct {
	err error
}

// returns a tea.Cmd that opens url in the browser and reports
// back whether it succeeded.
func openURLCmd(url string) tea.Cmd {
	return func() tea.Msg {
		return browserResultMsg{err: openURL(url)}
	}
}
