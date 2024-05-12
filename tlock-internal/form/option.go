package form

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tlockstyles "github.com/eklairs/tlock/tlock/styles"
)

// Keybinding for moving the choosen item to right
var KEY_RIGHT = key.NewBinding(
	key.WithKeys("right"),
	key.WithHelp("→", "right"),
)

// Keybinding for moving the choosen item to right
var KEY_LEFT = key.NewBinding(
	key.WithKeys("left"),
	key.WithHelp("←", "left"),
)

// Form item for option box
type FormItemOptionBox struct {
	// Title
	Title string

	// Description
	Description string

	// Input box
	Values []string

	// Selected value
	SelectedIndex int

	// Error message
	Focused bool
}

// Value
func (item FormItemOptionBox) Value() string {
	return item.Values[item.SelectedIndex]
}

// SetError
func (item FormItemOptionBox) SetError(err *error) {}

// Update
func (item *FormItemOptionBox) Update(msg tea.Msg) tea.Cmd {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, KEY_RIGHT):
			if item.SelectedIndex != len(item.Values)-1 {
				item.SelectedIndex += 1
			}

		case key.Matches(msgType, KEY_LEFT):
			if item.SelectedIndex != 0 {
				item.SelectedIndex -= 1
			}
		}
	}

	return nil
}

// Focus
func (item *FormItemOptionBox) Focus() {
	item.Focused = true
}

// Unfocus
func (item *FormItemOptionBox) Unfocus() {
	item.Focused = false
}

// View
func (item FormItemOptionBox) View() string {
	renderables := make([]string, len(item.Values))

	for index, option := range item.Values {
		if index == item.SelectedIndex {
			if item.Focused {
				renderables[index] = tlockstyles.Styles.AccentBgItem.Render(option)
			} else {
				renderables[index] = tlockstyles.Styles.SubAltBg.Render(option)
			}
		} else {
			renderables[index] = tlockstyles.Styles.SubTextItem.Render(option)
		}

		renderables[index] += " "
	}

	ui := lipgloss.JoinVertical(
		lipgloss.Left,
		tlockstyles.Styles.Title.Render(item.Title),
		tlockstyles.Styles.SubText.Render(item.Description), "",
		lipgloss.JoinHorizontal(lipgloss.Left, renderables...),
	)

	return lipgloss.NewStyle().Width(31).Render(ui)
}
