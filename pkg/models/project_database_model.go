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

type ProjectDatabase struct {
	inputs  map[string]textinput.Model
	focused string
	err     error
	cursor  int
	help    help.Model
	eggl    *EggLog
}

func ProjectDatabaseModel(l *EggLog) ProjectDatabase {
	return ProjectDatabase{
		inputs: map[string]textinput.Model{
			state.DatabaseURLName:      DatabaseURLInput(),
			state.DatabaseSqlcOrGoName: DatabaseSqlOrGoSQLCInput(),
			state.DatabaseRootName:     DatabaseRootLocation(),
		},
		focused: state.DatabaseURLName,
		err:     nil,
		help:    help.New(),
		cursor:  0,
		eggl:    l,
	}
}

func (m ProjectDatabase) IsUnsavedChanges() bool {
	return m.inputs[state.DatabaseURLName].Value() != state.DatabaseURL ||
		m.inputs[state.DatabaseSqlcOrGoName].Value() != state.DatabaseSqlcOrGo ||
		m.inputs[state.DatabaseRootName].Value() != state.DatabaseRoot
}

func (m ProjectDatabase) FocusFirstInput() {
	first := m.inputs[state.DatabaseURLName]
	first.Focus()
	m.focused = state.DatabaseURLName
	m.cursor = 0
	m.inputs[state.DatabaseURLName] = first
}

func (m ProjectDatabase) Init() tea.Cmd {
	return nil
}

func (m ProjectDatabase) View() string {
	return fmt.Sprintf(
		`%s %s

%s --> %s
%s --> %s
%s --> %s

`,

		styles.Keyword.Width(70).Render(" Postgres Settings"),
		NewUnsavedChangesIcon(m).View(),
		styles.Keyword.Align(lipgloss.Left).Width(30).Render("Database URL"),
		m.inputs[state.DatabaseURLName].View(),
		styles.Keyword.Align(lipgloss.Left).Width(30).Render("Write Go/SQL for SQLC"),
		m.inputs[state.DatabaseSqlcOrGoName].View(),
		styles.Keyword.Align(lipgloss.Left).Width(30).Render("Directory for Database Logic"),
		m.inputs[state.DatabaseRootName].View(),
	)
}

// nextInput focuses the next input field
func (m *ProjectDatabase) nextInput() {
	state_0 := state.DatabaseMap[m.cursor]
	m.cursor = (m.cursor + 1) % len(m.inputs)
	m.focused = state.DatabaseMap[m.cursor]
	previous := m.inputs[state_0]
	previous.Blur()
	next := m.inputs[m.focused]
	next.Focus()
	m.inputs[state_0] = previous
	m.inputs[m.focused] = next
}

// prevInput focuses the previous input field
func (m *ProjectDatabase) prevInput() {
	state_0 := state.DatabaseMap[m.cursor]
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.inputs) - 1
	}
	m.focused = state.DatabaseMap[m.cursor]
	previous := m.inputs[state_0]
	previous.Blur()
	next := m.inputs[m.focused]
	next.Focus()
	m.inputs[state_0] = previous
	m.inputs[m.focused] = next
}

func (m ProjectDatabase) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			state.DatabaseURL = m.inputs[state.DatabaseURLName].Value()
			state.DatabaseSqlcOrGo = m.inputs[state.DatabaseSqlcOrGoName].Value()
			state.DatabaseRoot = m.inputs[state.DatabaseRootName].Value()
		case tea.KeyEnter:
			if m.cursor == len(m.inputs)-1 {
				state.DatabaseURL = m.inputs[state.DatabaseURLName].Value()
				state.DatabaseSqlcOrGo = m.inputs[state.DatabaseSqlcOrGoName].Value()
				state.DatabaseRoot = m.inputs[state.DatabaseRootName].Value()
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
	for i := range state.DatabaseMap {
		m.inputs[state.DatabaseMap[i]], cmds[i] = m.inputs[state.DatabaseMap[i]].Update(msg)
	}
	return m, tea.Batch(cmds...)
}
