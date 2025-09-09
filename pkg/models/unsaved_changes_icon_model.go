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

// generate a model with bubbletea.Model
// to create a one character icon to be used if the 
// SubModel has unsaved changes that are different from 
// the tarrget state

import (
	tea "github.com/charmbracelet/bubbletea"
)

type UnsavedChangesIcon struct {
	// TODO
	Parent ISubModel
	UnsavedChanges bool
}

func NewUnsavedChangesIcon(parent ISubModel) UnsavedChangesIcon {
	return UnsavedChangesIcon{
		Parent: parent,
		UnsavedChanges: false,
	}
}

func (u UnsavedChangesIcon) Init() tea.Cmd {
	// TODO
	return nil
}

func (u UnsavedChangesIcon) View() string {
	// TODO
	if u.UnsavedChanges {
		return "󰆔"
	}
	return " "
}

func (u UnsavedChangesIcon) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO
	if u.Parent.IsUnsavedChanges() {
		u.UnsavedChanges = true
	}
	return u, nil 
}
