package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

// Creates a new input box
func InitializeInputBox(placeholder string) textinput.Model {
	input := textinput.New()
	input.Prompt = ""
	input.Width = 58
	input.Placeholder = placeholder
	input.PlaceholderStyle = tlockstyles.Styles.Placeholder

	return input
}
