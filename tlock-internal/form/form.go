package form

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/utils"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

// Form item
type FormItem interface {
	// Handle messages
	Update(msg tea.Msg) (FormItem, tea.Cmd)

	// View
	View() string

	// Focus
	Focus() FormItem

	// Unfocus
	Unfocus() FormItem

	// Returns the value
	Value() string
}

// ==== Form Item Input Box Start ====

// Form item for input boxes
type FormItemInputBox struct {
	// Title
	Title string

	// Description
	Description string

	// Input box
	Input textinput.Model

	// Error message
	ErrorMessage *string
}

// Update
func (item FormItemInputBox) Update(msg tea.Msg) (FormItem, tea.Cmd) {
	var cmd tea.Cmd

	item.Input, cmd = item.Input.Update(msg)

	return item, cmd
}

// Focus
func (item FormItemInputBox) Focus() FormItem {
	item.Input.Focus()

	return item
}

// Unfocus
func (item FormItemInputBox) Unfocus() FormItem {
	item.Input.Blur()

	return item
}

// View
func (item FormItemInputBox) View() string {
	return components.InputGroup(item.Title, item.Description, item.ErrorMessage, item.Input)
}

// Value
func (item FormItemInputBox) Value() string {
	return item.Input.Value()
}

// ==== Form Item Input Box End ====

// ==== Form Item Option Box Start ====

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

func (item FormItemOptionBox) Value() string {
	return item.Values[item.SelectedIndex]
}

// Update
func (item FormItemOptionBox) Update(msg tea.Msg) (FormItem, tea.Cmd) {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "right":
			if item.SelectedIndex != len(item.Values)-1 {
				item.SelectedIndex += 1
			}

		case "left":
			if item.SelectedIndex != 0 {
				item.SelectedIndex -= 1
			}
		}
	}

	return item, nil
}

// Focus
func (item FormItemOptionBox) Focus() FormItem {
	item.Focused = true

	return item
}

// Unfocus
func (item FormItemOptionBox) Unfocus() FormItem {
	item.Focused = false

	return item
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

// ==== Form Item Option Box End ====

type FormItemWrapped struct {
	// Form item
	FormItem FormItem

	// Is it enabled
	Enabled bool
}

type Form struct {
	// Items
	Items []FormItemWrapped

	// Focused index
	FocusedIndex int
}

// Initializes a new form item
func New(items []FormItem) Form {
	// Map
	wrapped := utils.Map(items, func(item FormItem) FormItemWrapped { return FormItemWrapped{FormItem: item, Enabled: true} })

	// Return instance
	form := Form{
		Items: wrapped,
	}

	// Focus first item
	form.Items[0].FormItem = form.Items[0].FormItem.Focus()

	// Return
	return form
}

// Switches focus from one index to another
func (form *Form) switchFocus(old, new int) {
	form.Items[old].FormItem = form.Items[old].FormItem.Unfocus()
	form.Items[new].FormItem = form.Items[new].FormItem.Focus()

	form.FocusedIndex = new
}

// Updates
func (form *Form) Update(msg tea.Msg) {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "tab":
			if form.FocusedIndex != len(form.Items)-1 {
				next := form.FocusedIndex + 1

				// If the next is disabled, then switch to its next
				if !form.Items[next].Enabled {
					next += 1
				}

				// Change focus
				form.switchFocus(form.FocusedIndex, next)
			}
		case "shift+tab":
			if form.FocusedIndex != 0 {
				next := form.FocusedIndex - 1

				// If the next is disabled, then switch to its next
				if !form.Items[next].Enabled && 0 != next-1 {
					next -= 1
				}

				// Change focus
				form.switchFocus(form.FocusedIndex, next)
			}
		default:
			form.Items[form.FocusedIndex].FormItem, _ = form.Items[form.FocusedIndex].FormItem.Update(msg)
		}
	}
}
