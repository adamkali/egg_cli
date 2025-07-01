package styles

import "github.com/charmbracelet/lipgloss"

var Hint = lipgloss.NewStyle().
	Bold(false).
	Underline(false).
	Italic(true).
	AlignHorizontal(lipgloss.Center).
	Foreground(HintColor)
