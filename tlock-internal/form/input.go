package form

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/eklairs/tlock/tlock-internal/components"
)

// Form item for input boxes
type FormItemInputBox struct {
	// Title
	Title string

	// Description
	Description string

	// Input box
	Input textinput.Model

	// Error message
	ErrorMessage *error
}

// Update
func (item *FormItemInputBox) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	// Let the input box handle its logic
	item.Input, cmd = item.Input.Update(msg)

	// Return
	return cmd
}

// Focus
func (item *FormItemInputBox) Focus() {
	item.Input.Focus()
}

// Unfocus
func (item *FormItemInputBox) Unfocus() {
	item.Input.Blur()
}

// View
func (item FormItemInputBox) View() string {
	return components.InputGroup(item.Title, item.Description, item.ErrorMessage, item.Input)
}

// SetError
func (item *FormItemInputBox) SetError(err *error) {
	item.ErrorMessage = err
}

// Value
func (item FormItemInputBox) Value() string {
	return item.Input.Value()
}
