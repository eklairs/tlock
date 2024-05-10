package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Message that represents that the form is submitted and the items pass the validators
type FormSubmittedMsg struct {
    Data map[string]string
}

// Validator
type Validator func(value string) error

// Form item wrapped
type FormItemWrapped struct {
    // ID
    ID string

	// Form item
	FormItem FormItem

	// Is it enabled
	Enabled bool

    // Validators
    Validators []Validator
}

type Form struct {
	// Items
	Items []FormItemWrapped

	// Focused index
	FocusedIndex int
}

// Initializes a new form item
func New(items []FormItem) Form {
	return Form{}
}

// Switches focus from one index to another
func (form *Form) switchFocus(old, new int) {
    // Do focus changing
	form.Items[old].FormItem.Unfocus()
	form.Items[new].FormItem.Focus()

    // Set new index
	form.FocusedIndex = new
}

// Updates
func (form *Form) Update(msg tea.Msg) tea.Cmd {
    cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		match_key: switch msgType.String() {
		case "tab":
			if form.FocusedIndex != len(form.Items)-1 && form.FocusedItem().ShouldGoNext() {
				next := form.FocusedIndex + 1

				// If the next is disabled, then switch to its next
				if !form.Items[next].Enabled {
					next += 1
				}

				// Change focus
				form.switchFocus(form.FocusedIndex, next)
			}
		case "shift+tab":
			if form.FocusedIndex != 0 && form.FocusedItem().ShouldGoPrev() {
				next := form.FocusedIndex - 1

				// If the next is disabled, then switch to its next
				if !form.Items[next].Enabled && 0 != next-1 {
					next -= 1
				}

				// Change focus
				form.switchFocus(form.FocusedIndex, next)
			}
        case "enter":
            data := make(map[string]string)

            // Validate them all!
            for _, item := range form.Items {
                // Run validators
                for _, validator := range item.Validators {
                    if err := validator(item.FormItem.Value()); err != nil {
                        // Set erro
                        item.FormItem.SetError(&err)

                        // There is issue with the form, break from the switch
                        break match_key
                    }
                }

                // Set the item
                data[item.ID] = item.FormItem.Value()
            }

            // Return
            cmds = append(cmds, func() tea.Msg { return FormSubmittedMsg{ Data: data } })
		}
	}

    // Let the current form item handle 
    cmds = append(cmds, form.Items[form.FocusedIndex].FormItem.Update(msg))

    // Return
    return tea.Batch(cmds...)
}

// Returns the focused form item
func (form *Form) FocusedItem() FormItem {
    return form.Items[form.FocusedIndex].FormItem
}

// Renders the form
func (form Form) View() string {
    // Items
    items := make([]string, len(form.Items))

    // Render them all
    for index, item := range form.Items {
        items[index] = item.FormItem.View()
    }

    // Render
    return lipgloss.JoinVertical(lipgloss.Center, items...)
}
