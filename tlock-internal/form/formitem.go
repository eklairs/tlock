package form

import tea "github.com/charmbracelet/bubbletea"

// Interface that every form item must implement
// Form item
type FormItem interface {
	// Handle messages
	Update(msg tea.Msg) tea.Cmd

	// View
	View() string

	// Focus
	Focus()

	// Unfocus
	Unfocus()

	// Returns the value
	Value() string

    // Sets the error message
    // Nil means to remove the error
    SetError(err *error)

    // Returns if the focus should go to the next item
    // Useful for groupping
    ShouldGoNext() bool

    // Returns if the focus should go to the next item
    // Useful for groupping
    ShouldGoPrev() bool
}
