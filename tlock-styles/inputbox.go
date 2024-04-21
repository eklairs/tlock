package tlockstyles

import (
	"github.com/charmbracelet/bubbles/textinput"
)

// Creates a new input box
func InitializeInputBox(styles Styles, placeholder string) textinput.Model {
    return InitializeInputBoxCustomWidth(styles, placeholder, 58)
}

// Creates a new input box
func InitializeInputBoxCustomWidth(styles Styles, placeholder string, width int) textinput.Model {
	input := textinput.New()
	input.Prompt = ""
	input.Width = width
	input.Placeholder = placeholder
	input.PlaceholderStyle = styles.InputPlaceholder

	return input
}

