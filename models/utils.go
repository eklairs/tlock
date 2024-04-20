package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// Creates a new input box
func InitializeInputBox(placeholder string) textinput.Model {
    input := textinput.New();
    input.Prompt = ""
    input.Width = 58
    input.Placeholder = placeholder
    input.PlaceholderStyle = lipgloss.NewStyle().Background(COLOR_BG_OVER).Foreground(COLOR_DIMMED)

    return input
}

