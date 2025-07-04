package models

import (
	"fmt"

	"github.com/adamkali/egg_cli/state"
	"github.com/adamkali/egg_cli/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ProjectSettings struct {
	inputs  map[string]textinput.Model
	focused string
	err     error
	cursor  int
	help    help.Model
	eggl    *EggLog
}

func ProjectSettingsModel(l *EggLog) ProjectSettings {
	return ProjectSettings{
		inputs: map[string]textinput.Model{
			state.ProjectHostName:     ProjectHostInput(),
			state.ProjectUsernameName: ProjectUsernameInput(),
			state.ProjectNameName:     ProjectNameInput(),
		},
		focused: state.ProjectHostName,
		err:     nil,
		help:    help.New(),
		cursor:  0,
		eggl:    l,
	}
}

func (m ProjectSettings) FocusFirstInput() {
	first := m.inputs[state.ProjectHostName]
	first.Focus()
	m.focused = state.ProjectHostName
	m.cursor = 0
	m.inputs[state.ProjectHostName] = first
}

func (m ProjectSettings) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m ProjectSettings) View() string {
	return fmt.Sprintf(
		`%s %s

%s/%s/%s
%s/%s/%s


`,
		styles.Keyword.Width(70).Render(" Project Settings"),
		NewUnsavedChangesIcon(m).View(),
		styles.Keyword.Render("hosting"),
		styles.Keyword.Render("username"),
		styles.Keyword.Render("project"),
		m.inputs[state.ProjectHostName].View(),
		m.inputs[state.ProjectUsernameName].View(),
		m.inputs[state.ProjectNameName].View(),
	)
}

// nextInput focuses the next input field
func (m *ProjectSettings) nextInput() {
	state_0 := state.ProjectSettingsMap[m.cursor]
	m.cursor = (m.cursor + 1) % len(m.inputs)
	m.focused = state.ProjectSettingsMap[m.cursor]
	m.eggl.Info("nextInput: cursor=%d, focused=%s", m.cursor, m.focused)
	previous := m.inputs[state_0]
	previous.Blur()
	next := m.inputs[m.focused]
	next.Focus()
	m.inputs[state_0] = previous
	m.inputs[m.focused] = next
}

// prevInput focuses the previous input field
func (m *ProjectSettings) prevInput() {
	state_0 := state.ProjectSettingsMap[m.cursor]
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.inputs) - 1
	}
	m.focused = state.ProjectSettingsMap[m.cursor]
	m.eggl.Info("prevInput: cursor=%d, focused=%s", m.cursor, m.focused)
	previous := m.inputs[state_0]
	previous.Blur()
	next := m.inputs[m.focused]
	next.Focus()
	m.inputs[state_0] = previous
	m.inputs[m.focused] = next
}

func (m ProjectSettings) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			state.ProjectHost = m.inputs[state.ProjectHostName].Value()
			state.ProjectUsername = m.inputs[state.ProjectUsernameName].Value()
			state.ProjectName = m.inputs[state.ProjectNameName].Value()
			state.ProjectNamespace = state.ProjectHost + "/" +
				state.ProjectUsername + "/" +
				state.ProjectName
		case tea.KeyEnter:
			if m.cursor == len(m.inputs)-1 {
				state.ProjectHost = m.inputs[state.ProjectHostName].Value()
				state.ProjectUsername = m.inputs[state.ProjectUsernameName].Value()
				state.ProjectName = m.inputs[state.ProjectNameName].Value()
				state.ProjectNamespace = state.ProjectHost + "/" +
					state.ProjectUsername + "/" +
					state.ProjectName
			}
			m.nextInput()
		case tea.KeyCtrlC:
			return m, tea.Quit
		// probably add this to the form as well
		case tea.KeyShiftTab:
			m.prevInput()
		case tea.KeyTab:
			m.nextInput()
		}
	// We handle errors just like any other message
	case state.ErrMsg:
		m.err = msg
		return m, nil
	}

	for i := range state.ProjectSettingsMap {
		m.inputs[state.ProjectSettingsMap[i]], cmds[i] = m.inputs[state.ProjectSettingsMap[i]].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m ProjectSettings) IsUnsavedChanges() bool {
	return m.inputs[state.ProjectHostName].Value() != state.ProjectHost ||
		m.inputs[state.ProjectUsernameName].Value() != state.ProjectUsername ||
		m.inputs[state.ProjectNameName].Value() != state.ProjectName
}
