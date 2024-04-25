package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

func InputGroup(title, description string, error *string, input textinput.Model) string {
	// Total width relative to the input's width
	width := input.Width + 7

	items := []string{
		tlockstyles.Styles.Title.Copy().Width(width).Render(title),
		tlockstyles.Styles.SubText.Copy().Width(width).Render(description),
		tlockstyles.Styles.Input.Copy().Width(width).Render(input.View()), "",
	}

	// Append error if any
	if error != nil {
		items = append(items, tlockstyles.Styles.Error.Copy().Width(65).Render(fmt.Sprintf("× %s", *error)), "")
	}

	return lipgloss.JoinVertical(lipgloss.Center, items...)
}
