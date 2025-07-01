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
