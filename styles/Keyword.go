package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var Keyword = lipgloss.NewStyle().
	Bold(true).
	AlignHorizontal(lipgloss.Center).
	Foreground(lipgloss.Color("#ff8933")).
	PaddingLeft(1)
