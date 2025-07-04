package models

import (
	"fmt"

	"github.com/adamkali/egg_cli/state"
	"github.com/adamkali/egg_cli/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ProjectServerSettings struct {
	inputs  map[string]textinput.Model
	focused string
	err     error
	cursor  int
	help    help.Model
	eggl    *EggLog
}

func ProjectServerSettingsModel(log *EggLog) ProjectServerSettings {
	return ProjectServerSettings{
		inputs: map[string]textinput.Model{
			state.ServerJWTName:       ServerJWTInput(),
			state.ServerPortName:      ServerPortInput(),
		},
		focused: state.ServerJWTName,
		err:     nil,
		help:    help.New(),
		cursor:  0,
		eggl:    log,
	}
}

func (m ProjectServerSettings) Init() tea.Cmd {
	return nil
}

func (m ProjectServerSettings) View() string {
	return fmt.Sprintf(
		`%s %s

%s --> %s
%s --> %s
%s --> %s
%s --> %s
`,
		styles.Keyword.Render("󰒋 Server Settings"),
		NewUnsavedChangesIcon(m).View(),
		styles.Keyword.Width(30).Render("JWT"),
		m.inputs[state.ServerJWTName].View(),
		styles.Keyword.Width(30).Render("Server Port"),
		m.inputs[state.ServerPortName].View(),
		styles.Keyword.Width(30).Render("Frontend Dist Directory"),
		m.inputs[state.ServerFrontendDirName].View(),
		styles.Keyword.Width(30).Render("Frontend Api Directory"),
		m.inputs[state.ServerFrontendApiName].View(),
	)
}

func (m ProjectServerSettings) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			state.ServerJWT = m.inputs[state.ServerJWTName].Value()
			state.ServerPort = m.inputs[state.ServerPortName].Value()
			state.ServerFrontendDir = m.inputs[state.ServerFrontendDirName].Value()
			state.ServerFrontendApi = m.inputs[state.ServerFrontendApiName].Value()
		case tea.KeyEnter:
			m.nextInput()		
		case tea.KeyCtrlC:
			return m, tea.Quit	
		case tea.KeyShiftTab:
			m.prevInput()
		case tea.KeyTab:
			m.nextInput()
		}
	case state.ErrMsg:
		m.err = msg
		return m, nil
	}

	m.eggl.Info("len(state.ProjectSeverMap): %d", len(state.ProjectSeverMap))
	for i := range state.ProjectSeverMap {
		m.inputs[state.ProjectSeverMap[i]], cmds[i] = m.inputs[state.ProjectSeverMap[i]].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m ProjectServerSettings) IsUnsavedChanges() bool {
	return m.inputs[state.ServerJWTName].Value() != state.ServerJWT ||
		m.inputs[state.ServerPortName].Value() != state.ServerPort ||
		m.inputs[state.ServerFrontendDirName].Value() != state.ServerFrontendDir ||
		m.inputs[state.ServerFrontendApiName].Value() != state.ServerFrontendApi
}

func (m ProjectServerSettings) FocusFirstInput() {
	first := m.inputs[state.ServerJWTName]
	first.Focus()
	m.focused = state.ServerJWTName // TODO: set this.
	m.cursor = 0
	m.inputs[state.ServerJWTName] = first
}


func (m *ProjectServerSettings) nextInput() {
	state_0 := state.ProjectSeverMap[m.cursor]
	m.cursor = (m.cursor + 1) % len(m.inputs)
	m.focused = state.ProjectSeverMap[m.cursor]
	m.eggl.Info("nextInput: cursor=%d, focused=%s", m.cursor, m.focused)
	previous := m.inputs[state_0]
	previous.Blur()
	next := m.inputs[m.focused]
	next.Focus()
	m.inputs[state_0] = previous
	m.inputs[m.focused] = next
}

// prevInput focuses the previous input field
func (m *ProjectServerSettings) prevInput() {
	state_0 := state.ProjectSeverMap[m.cursor]
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.inputs) - 1
	}
	m.focused = state.ProjectSeverMap[m.cursor]
	previous := m.inputs[state_0]
	previous.Blur()
	next := m.inputs[m.focused]
	next.Focus()
	m.inputs[state_0] = previous
	m.inputs[m.focused] = next
}
