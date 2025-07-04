package models

import (
	"fmt"
	"strings"

	"github.com/adamkali/egg_cli/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PageModel struct {
	pages       []tea.Model
	currentPage int
	eggl        *EggLog
}

func CreatePageModel(log *EggLog) PageModel {
	settings := ProjectSettingsModel(log)
	settings.FocusFirstInput()
	database := ProjectDatabaseModel(log)
	license := ProjectLicenseModel(log)
	server := ProjectServerSettingsModel(log)
	mini := ProjectS3Model(log)

	return PageModel{
		pages: []tea.Model{
			settings,
			license,
			server,
			database,
			mini,
		},
		currentPage: 0,
		eggl:        log,
	}
}

func (m PageModel) Init() tea.Cmd {
	m.eggl.Info("PageModel started")
	return nil
}

func (m PageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+n":
			// Move to next page
			if m.currentPage <= len(m.pages)-1 {
				m.currentPage++
				if m.currentPage > len(m.pages)-1 {
					return m, tea.Quit
				}
				m.eggl.Info("Moving to page %d", m.currentPage)
				if page, ok := m.pages[m.currentPage].(ISubModel); ok {
					page.FocusFirstInput()
				}
			}
		case "ctrl+p":
			// Move to previous page
			if m.currentPage > 0 {
				m.currentPage--
				m.eggl.Info("Moving to page %d", m.currentPage)
				// Focus first input of new page
				if page, ok := m.pages[m.currentPage].(ISubModel); ok {
					page.FocusFirstInput()
				}
			}
		}
	}

	// Update current page
	updatedPage, cmd := m.pages[m.currentPage].Update(msg)
	m.pages[m.currentPage] = updatedPage

	return m, cmd
}

func (m PageModel) View() string {
	styleModal := lipgloss.NewStyle().
		BorderBottom(true).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderStyle(lipgloss.RoundedBorder()).
		PaddingBottom(1)

	view := fmt.Sprintf("%s", styles.TitleStyle.Width(100).Render("🥚 Initialize a Project"))
	view += "\n" + m.headerLine() + "\n"
	view += styleModal.Render(m.pages[m.currentPage].View()) + "\n"
	view += m.footerLine()

	return view
}

func (m *PageModel) headerLine() string {
	pageNames := []string{
		"Settings",
		"Server",
		"License",
		"Database",
		"S3",
	}
	tabs := make([]string, 0)

	for i, name := range pageNames {
		if i == m.currentPage {
			style := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#ff8933"))

			name = style.Render(fmt.Sprintf(" %s ", name))
		}

		tabs = append(tabs, lipgloss.NewStyle().Padding(0, 1).Render(name))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

func (m *PageModel) footerLine() string {
	styleItem := lipgloss.
		NewStyle().
		Width(18).
		Align(lipgloss.Center)
	styleLabel := lipgloss.
		NewStyle().
		Padding(0, 1).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#ff8933")).
		Align(lipgloss.Right)
	styleHotkey := lipgloss.
		NewStyle().
		Padding(0, 1).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#ff8933")).
		Bold(true).
		Align(lipgloss.Left)

	var footer string

	hotkeys := []string{
		"Previous::ctrl+p",
		"Next::ctrl+n",
		"Save::ctrl+s",
		"Exit::ctrl+c",
		"Help::ctrl+h",
		"+ Input::Tab",
		"- Input::Shift+Tab",
		"Quit::ctrl+q",
	}

	for _, entry := range hotkeys {
		label := strings.Split(entry, "::")[0]
		hotkey := strings.Split(entry, "::")[1]
		label = styleLabel.Render(label)
		hotkey = styleHotkey.Render(hotkey)
		footer += styleItem.Render(label + hotkey)
	}

	return footer
}
