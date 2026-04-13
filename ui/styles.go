package ui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("63"))

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	modeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("214"))

	counterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	counterOverStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("196"))

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	recentHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("39"))

	recentItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	statusErrStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	statusOkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46"))

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1)
)
