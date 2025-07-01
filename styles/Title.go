package styles
import "github.com/charmbracelet/lipgloss"
var TitleStyle = lipgloss.NewStyle().
	Bold(true).
	AlignHorizontal(lipgloss.Center).
	Foreground(lipgloss.Color("#ff8933")).
	PaddingTop(1).
	PaddingLeft(4).
	PaddingBottom(1).
	PaddingRight(4)

