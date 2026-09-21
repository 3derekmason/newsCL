package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"newscl/internal/hn"
)

const (
	modalOptionLink = iota
	modalOptionComments
)

// state for the item detail popup
type modalModel struct {
	item      hn.Item
	cursor    int
	status    string
	statusErr bool
}

func newModalModel(item hn.Item) modalModel {
	return modalModel{item: item}
}

// opens the detail popup for item on top of the current list
func (m Model) enterModal(item hn.Item) Model {
	m.screen = screenModal
	m.modal = newModalModel(item)
	return m
}

func (m Model) updateModal(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case browserResultMsg:
		if msg.err != nil {
			m.modal.status = "Couldn't open browser: " + msg.err.Error()
			m.modal.statusErr = true
		} else {
			m.modal.status = "Opened in your browser."
			m.modal.statusErr = false
		}
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "left", "q":
			m.screen = screenList
		case "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.modal.cursor > 0 {
				m.modal.cursor--
			}
		case "down", "j":
			if m.modal.cursor < modalOptionComments {
				m.modal.cursor++
			}
		case "enter":
			m.modal.status = ""
			m.modal.statusErr = false
			var url string
			if m.modal.cursor == modalOptionLink {
				url = m.modal.item.LinkURL()
			} else {
				url = m.modal.item.DiscussionURL()
			}
			return m, openURLCmd(url)
		}
	}

	return m, nil
}

// composes the detail popup on top of the list screen behind it,
// using Lip Gloss v2's canvas/layer compositing
func (m Model) viewModal() string {
	background := m.viewList()
	dialog := m.viewModalDialog()

	width, height := m.width, m.height
	if width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 24
	}

	canvas := lipgloss.NewCanvas(width, height)
	canvas.Compose(lipgloss.NewLayer(background).X(0).Y(0))

	dialogWidth := lipgloss.Width(dialog)
	dialogHeight := lipgloss.Height(dialog)
	x := max(0, (width-dialogWidth)/2)
	y := max(0, (height-dialogHeight)/2)
	canvas.Compose(lipgloss.NewLayer(dialog).X(x).Y(y).Z(1))

	return canvas.Render()
}

func (m Model) viewModalDialog() string {
	item := m.modal.item

	linkLabel := "View Post"
	if item.IsJob() {
		linkLabel = "View Job"
	}

	options := []string{
		linkLabel,
		fmt.Sprintf("Comments (%d)", item.Descendants),
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render(truncate(item.Title, 64)))
	if d := domain(item.URL); d != "" {
		b.WriteString("\n" + mutedStyle.Render(d))
	}
	b.WriteString("\n\n")

	for i, opt := range options {
		b.WriteString(renderChoice(opt, i == m.modal.cursor))
		b.WriteString("\n")
	}

	if m.modal.status != "" {
		style := successStyle
		if m.modal.statusErr {
			style = errorStyle
		}
		b.WriteString("\n" + style.Render(m.modal.status) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("[↑/↓] select   [enter] open   [←] back"))

	return dialogStyle.Width(80).Render(b.String())
}
