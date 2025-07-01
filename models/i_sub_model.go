package models

type ISubModel interface {
	FocusFirstInput()
	IsUnsavedChanges() bool
}
