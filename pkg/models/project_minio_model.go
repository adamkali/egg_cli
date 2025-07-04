package models

import (
	"fmt"

	"github.com/adamkali/egg_cli/state"
	"github.com/adamkali/egg_cli/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ProjectS3 struct {
	inputs  map[string]textinput.Model
	focused string
	err     error
	cursor  int
	help    help.Model
	eggl    *EggLog
}

func ProjectS3Model(l *EggLog) ProjectS3 {
	return ProjectS3{
		inputs: map[string]textinput.Model{
			state.MinioURLName:   MinioUrlInput(),
			state.MinioAccessKeyName: MinioAccessKeyInput(),
			state.MinioSecretKeyName: MinioSecretKeyInput(),
		},
		focused: state.MinioURLName,
		err:     nil,
		help:    help.New(),
		cursor:  0,
		eggl:    l,
	}
}

func (m ProjectS3) IsUnsavedChanges() bool {
	return m.inputs[state.MinioURLName].Value() != state.MinioURL ||
		m.inputs[state.MinioAccessKeyName].Value() != state.MinioAccessKey ||
		m.inputs[state.MinioSecretKeyName].Value() != state.MinioSecretKey
}

func (m ProjectS3) FocusFirstInput() {
	first := m.inputs[state.MinioURLName]
	first.Focus()
	m.focused = state.MinioURLName
	m.cursor = 0
	m.inputs[state.MinioURLName] = first
}

func (m ProjectS3) Init() tea.Cmd {
	return nil
}

func (m ProjectS3) View() string {
	return fmt.Sprintf(
		`%s %s

%s --> %s
%s --> %s
%s --> %s

`,

		styles.Keyword.Width(70).Render("S3 Settings"), 
		NewUnsavedChangesIcon(m).View(),
		styles.Keyword.Width(30).Render("S3 URL"),
		m.inputs[state.MinioURLName].View(),
		styles.Keyword.Width(30).Render("Access Key"),
		m.inputs[state.MinioAccessKeyName].View(),
		styles.Keyword.Width(30).Render("Secret Key"),
		m.inputs[state.MinioSecretKeyName].View(),
	)
}

// nextInput focuses the next input field
func (m *ProjectS3) nextInput() {
	state_0 := state.MinioMap[m.cursor]
	m.cursor = (m.cursor + 1) % len(m.inputs)
	m.focused = state.MinioMap[m.cursor]
	previous := m.inputs[state_0]
	previous.Blur()
	next := m.inputs[m.focused]
	next.Focus()
	m.inputs[state_0] = previous
	m.inputs[m.focused] = next
}

// prevInput focuses the previous input field
func (m *ProjectS3) prevInput() {
	state_0 := state.MinioMap[m.cursor]
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.inputs) - 1
	}
	m.focused = state.MinioMap[m.cursor]
	m.eggl.Info("prevInput: cursor=%d, focused=%s", m.cursor, m.focused)
	previous := m.inputs[state_0]
	previous.Blur()
	next := m.inputs[m.focused]
	next.Focus()
	m.inputs[state_0] = previous
	m.inputs[m.focused] = next
}

func (m ProjectS3) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			state.MinioURL = m.inputs[state.MinioURLName].Value()
			state.MinioAccessKey = m.inputs[state.MinioAccessKeyName].Value()
			state.MinioSecretKey = m.inputs[state.MinioSecretKeyName].Value()
		case tea.KeyEnter:
			if m.cursor == len(m.inputs)-1 {
				state.MinioURL = m.inputs[state.MinioURLName].Value()
				state.MinioAccessKey = m.inputs[state.MinioAccessKeyName].Value()
				state.MinioSecretKey = m.inputs[state.MinioSecretKeyName].Value()
			}
			m.nextInput()
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyShiftTab:
			m.prevInput()
		case tea.KeyTab:
			m.nextInput()
		}

	// We handle errors just like any other message
	case state.ErrMsg:
		m.err = msg
		m.eggl.Error("error: %v", m.err)
		return m, nil
	}
	for i := range state.MinioMap {
		m.inputs[state.MinioMap[i]], cmds[i] = m.inputs[state.MinioMap[i]].Update(msg)
	}
	return m, tea.Batch(cmds...)
}
