/*
Copyright © 2025 Adam Kalinowski <adam.kalilarosa@proton.me>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package models

import (
	"fmt"

	"github.com/adamkali/egg_cli/state"
	"github.com/adamkali/egg_cli/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ProjectAuthModel struct {
	inputs  map[string]textinput.Model
	focused string
	err     error
	cursor  int
	help    help.Model
	eggl    *EggLog
}

func ProjectAuthSettingsModel(log *EggLog) ProjectAuthModel {
	return ProjectAuthModel{
		inputs: map[string]textinput.Model{
			state.AuthProviderName:      AuthProviderInput(),
			state.ServerJWTName:         ServerJWTInput(),
			state.Auth0DomainName:       Auth0DomainInput(),
			state.Auth0AudienceName:     Auth0AudienceInput(),
			state.Auth0ClientIDName:     Auth0ClientIDInput(),
			state.Auth0ClientSecretName: Auth0ClientSecretInput(),
		},
		focused: state.AuthProviderName,
		err:     nil,
		help:    help.New(),
		cursor:  0,
		eggl:    log,
	}
}

func (m ProjectAuthModel) Init() tea.Cmd { return nil }

func (m ProjectAuthModel) View() string {
	return fmt.Sprintf(
		`%s %s

%s --> %s
%s --> %s
%s --> %s
%s --> %s
%s --> %s
%s --> %s
`,
		styles.Keyword.Render("󰒃 Auth Settings"),
		NewUnsavedChangesIcon(m).View(),
		styles.Keyword.Width(30).Render("Auth Provider (jwt|auth0)"),
		m.inputs[state.AuthProviderName].View(),
		styles.Keyword.Width(30).Render("JWT Secret"),
		m.inputs[state.ServerJWTName].View(),
		styles.Keyword.Width(30).Render("Auth0 Domain"),
		m.inputs[state.Auth0DomainName].View(),
		styles.Keyword.Width(30).Render("Auth0 Audience"),
		m.inputs[state.Auth0AudienceName].View(),
		styles.Keyword.Width(30).Render("Auth0 Client ID"),
		m.inputs[state.Auth0ClientIDName].View(),
		styles.Keyword.Width(30).Render("Auth0 Client Secret"),
		m.inputs[state.Auth0ClientSecretName].View(),
	)
}

func (m ProjectAuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			state.AuthProvider = m.inputs[state.AuthProviderName].Value()
			if state.AuthProvider == "" {
				state.AuthProvider = "jwt"
			}
			state.ServerJWT = m.inputs[state.ServerJWTName].Value()
			state.Auth0Domain = m.inputs[state.Auth0DomainName].Value()
			state.Auth0Audience = m.inputs[state.Auth0AudienceName].Value()
			state.Auth0ClientID = m.inputs[state.Auth0ClientIDName].Value()
			state.Auth0ClientSecret = m.inputs[state.Auth0ClientSecretName].Value()
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

	for i := range state.AuthMap {
		m.inputs[state.AuthMap[i]], cmds[i] = m.inputs[state.AuthMap[i]].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m ProjectAuthModel) IsUnsavedChanges() bool {
	return m.inputs[state.AuthProviderName].Value() != state.AuthProvider ||
		m.inputs[state.ServerJWTName].Value() != state.ServerJWT ||
		m.inputs[state.Auth0DomainName].Value() != state.Auth0Domain ||
		m.inputs[state.Auth0AudienceName].Value() != state.Auth0Audience ||
		m.inputs[state.Auth0ClientIDName].Value() != state.Auth0ClientID ||
		m.inputs[state.Auth0ClientSecretName].Value() != state.Auth0ClientSecret
}

func (m ProjectAuthModel) FocusFirstInput() {
	first := m.inputs[state.AuthProviderName]
	first.Focus()
	m.inputs[state.AuthProviderName] = first
}

func (m *ProjectAuthModel) nextInput() {
	state0 := state.AuthMap[m.cursor]
	m.cursor = (m.cursor + 1) % len(m.inputs)
	m.focused = state.AuthMap[m.cursor]
	prev := m.inputs[state0]
	prev.Blur()
	next := m.inputs[m.focused]
	next.Focus()
	m.inputs[state0] = prev
	m.inputs[m.focused] = next
}

func (m *ProjectAuthModel) prevInput() {
	state0 := state.AuthMap[m.cursor]
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.inputs) - 1
	}
	m.focused = state.AuthMap[m.cursor]
	prev := m.inputs[state0]
	prev.Blur()
	next := m.inputs[m.focused]
	next.Focus()
	m.inputs[state0] = prev
	m.inputs[m.focused] = next
}
