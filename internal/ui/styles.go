package ui

import "charm.land/lipgloss/v2"

var (
	colorAccent  = lipgloss.Color("#f59e0b")
	colorMuted   = lipgloss.Color("#64748b")
	colorError   = lipgloss.Color("#dc2626")
	colorSuccess = lipgloss.Color("#10b981")
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	pointerStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(colorAccent).
				Bold(true)

	itemStyle = lipgloss.NewStyle()

	mutedStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	successStyle = lipgloss.NewStyle().
			Foreground(colorSuccess)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorError)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorAccent).
			Padding(1, 3)

	dialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSuccess).
			Padding(1, 3)
)
