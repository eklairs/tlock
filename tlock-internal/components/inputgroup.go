package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

func InputGroup(title, description string, error *string, input textinput.Model) string {
	items := []string{
		tlockstyles.Styles.Title.Copy().Width(65).Render(title),
		tlockstyles.Styles.SubText.Copy().Width(65).Render(description),
		tlockstyles.Styles.Input.Copy().Width(65).Render(input.View()), "",
	}

	// Append error if any
	if error != nil {
		items = append(items, tlockstyles.Styles.Error.Copy().Width(65).Render(fmt.Sprintf("× %s", *error)), "")
	}

	return lipgloss.JoinVertical(lipgloss.Center, items...)
}
