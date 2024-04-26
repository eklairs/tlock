package form

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tlockinternal "github.com/eklairs/tlock/tlock-internal"
	"github.com/eklairs/tlock/tlock-internal/components"
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
}

// ==== Form Item Input Box Start ====

// Form item for input boxes
type FormItemInputBox struct {
	// Title
	title string

	// Description
	description string

	// Input box
	input textinput.Model

	// Error message
	errorMessage *string

	// Width
	width int
}

// Update
func (item FormItemInputBox) Update(msg tea.Msg) (FormItem, tea.Cmd) {
	var cmd tea.Cmd

	item.input, cmd = item.input.Update(msg)

	return item, cmd
}

// Focus
func (item FormItemInputBox) Focus() FormItem {
	item.input.Focus()

	return item
}

// Unfocus
func (item FormItemInputBox) Unfocus() FormItem {
	item.input.Blur()

	return item
}

// View
func (item FormItemInputBox) View() string {
	return components.InputGroup(item.title, item.description, item.errorMessage, item.input)
}

// ==== Form Item Input Box End ====

// ==== Form Item Option Box Start ====

// Form item for option box
type FormItemOptionBox struct {
	// Title
	title string

	// Description
	description string

	// Input box
	values []string

	// Selected value
	selectedIndex int

	// Error message
	focused bool
}

func (item FormItemOptionBox) SelectedValue() string {
	return item.values[item.selectedIndex]
}

// Update
func (item FormItemOptionBox) Update(msg tea.Msg) (FormItem, tea.Cmd) {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "right":
			if item.selectedIndex != len(item.values)-1 {
				item.selectedIndex += 1
			}

		case "left":
			if item.selectedIndex != 0 {
				item.selectedIndex -= 1
			}
		}
	}

	return item, nil
}

// Focus
func (item FormItemOptionBox) Focus() FormItem {
	item.focused = true

	return item
}

// Unfocus
func (item FormItemOptionBox) Unfocus() FormItem {
	item.focused = false

	return item
}

// View
func (item FormItemOptionBox) View() string {
	renderables := make([]string, len(item.values))

	for index, option := range item.values {
		if index == item.selectedIndex {
			if item.focused {
				renderables[index] = tlockstyles.Styles.AccentBgItem.Render(option)
			} else {
				renderables[index] = tlockstyles.Styles.SubAltBg.Render(option)
			}
		} else {
			renderables[index] = tlockstyles.Styles.SubText.Render(option)
		}

		renderables[index] += " "
	}

	ui := lipgloss.JoinVertical(
		lipgloss.Left,
		tlockstyles.Styles.Title.Render(item.title),
		tlockstyles.Styles.SubText.Render(item.description), "",
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
	wrapped := tlockinternal.Map(items, func(item FormItem) FormItemWrapped { return FormItemWrapped{FormItem: item, Enabled: true} })

	// Return instance
	return Form{
		Items: wrapped,
	}
}

// Switches focus from one index to another
func (form *Form) switchFocus(old, new int) {
	form.Items[old].FormItem = form.Items[old].FormItem.Unfocus()
	form.Items[new].FormItem = form.Items[new].FormItem.Focus()

	form.FocusedIndex = new
}

// Updates
func (form *Form) Update(msg tea.Msg) {
	// Update the current focused item
	form.Items[form.FocusedIndex].FormItem, _ = form.Items[form.FocusedIndex].FormItem.Update(msg)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "tab":
			if form.FocusedIndex != len(form.Items)-1 {
				next := form.FocusedIndex + 1

				// If the next is disabled, then switch to its next
				if !form.Items[next].Enabled && len(form.Items)-1 != next+1 {
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
		}
	}
}
