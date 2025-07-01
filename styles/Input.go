package styles
import "github.com/charmbracelet/lipgloss"
var Input = lipgloss.NewStyle().
	Bold(false).
	Underline(true).
	AlignHorizontal(lipgloss.Center).
	Foreground(Foreground).
	PaddingLeft(1)
