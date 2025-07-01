package models

import (
	"fmt"

	"github.com/adamkali/egg_cli/state"
	"github.com/adamkali/egg_cli/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProjectLicense struct {
	inputs  map[string]textinput.Model
	focused string
	err     error
	cursor  int
	help    help.Model
	eggl    *EggLog
}

func ProjectLicenseModel(l *EggLog) ProjectLicense {
	return ProjectLicense{
		inputs: map[string]textinput.Model{
			state.LicenseName:         LicenseInput(),
			state.CopyrightYearName:   CopyrightYearInput(),
			state.CopyrightAuthorName: CopyrightAuthorInput(),
		},
		focused: state.LicenseName,
		err:     nil,
		help:    help.New(),
		cursor:  0,
		eggl:    l,
	}
}

func (m ProjectLicense) FocusFirstInput() {
	first := m.inputs[state.LicenseName]
	first.Focus()
	m.focused = state.LicenseName
	m.cursor = 0
	m.inputs[state.LicenseName] = first
}

func (m ProjectLicense) IsUnsavedChanges() bool {
	return m.inputs[state.LicenseName].Value() != state.License ||
		m.inputs[state.CopyrightYearName].Value() != state.CopyrightYear ||
		m.inputs[state.CopyrightAuthorName].Value() != state.CopyrightAuthor
}

func (m ProjectLicense) Init() tea.Cmd {
	return nil
}

func (m ProjectLicense) View() string {
	view := fmt.Sprintf(
		`%s %s

%s --> %s
%s --> %s
%s --> %s

`,

		styles.Keyword.Width(70).Render(" License & Copyright"),
		NewUnsavedChangesIcon(m).View(),
		styles.Keyword.Align(lipgloss.Left).Width(30).Render("License"),
		m.inputs[state.LicenseName].View(),
		styles.Keyword.Align(lipgloss.Left).Width(30).Render("Copyright Year"),
		m.inputs[state.CopyrightYearName].View(),
		styles.Keyword.Align(lipgloss.Left).Width(30).Render("Copyright Author"),
		m.inputs[state.CopyrightAuthorName].View(),
	)

	return view
}

// nextInput focuses the next input field
func (m *ProjectLicense) nextInput() {
	state_0 := state.LicenseMap[m.cursor]
	m.cursor = (m.cursor + 1) % len(m.inputs)
	m.focused = state.LicenseMap[m.cursor]
	m.eggl.Info("nextInput: cursor=%d, focused=%s", m.cursor, m.focused)
	previous := m.inputs[state_0]
	previous.Blur()
	next := m.inputs[m.focused]
	next.Focus()
	m.inputs[state_0] = previous
	m.inputs[m.focused] = next
}

// prevInput focuses the previous input field
func (m *ProjectLicense) prevInput() {
	state_0 := state.LicenseMap[m.cursor]
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.inputs) - 1
	}
	m.focused = state.LicenseMap[m.cursor]
	m.eggl.Info("prevInput: cursor=%d, focused=%s", m.cursor, m.focused)
	previous := m.inputs[state_0]
	previous.Blur()
	next := m.inputs[m.focused]
	next.Focus()
	m.inputs[state_0] = previous
	m.inputs[m.focused] = next
}

func (m ProjectLicense) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			state.License = m.inputs[state.LicenseName].Value()
			state.CopyrightYear = m.inputs[state.CopyrightYearName].Value()
			state.CopyrightAuthor = m.inputs[state.CopyrightAuthorName].Value()
		case tea.KeyEnter:
			if m.cursor == len(m.inputs)-1 {
				state.License = m.inputs[state.LicenseName].Value()
				state.CopyrightYear = m.inputs[state.CopyrightYearName].Value()
				state.CopyrightAuthor = m.inputs[state.CopyrightAuthorName].Value()
				// Move to next page will be handled by parent
			}
			m.nextInput()
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyShiftTab:
			if m.cursor == 0 {
				// Move to previous page will be handled by parent
			}
			m.prevInput()
		case tea.KeyTab:
			m.nextInput()
		}

	// We handle errors just like any other message
	case state.ErrMsg:
		m.err = msg
		return m, nil
	}

	// Clear error when user starts typing
	m.err = nil

	for i := range state.LicenseMap {
		updatedInput, cmd := m.inputs[state.LicenseMap[i]].Update(msg)
		m.inputs[state.LicenseMap[i]] = updatedInput
		cmds[i] = cmd
	}

	return m, tea.Batch(cmds...)
}
