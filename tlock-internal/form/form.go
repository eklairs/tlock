package form

import (
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

// Validator
type Validator = func(vault *tlockvault.Vault, value string) error

// Message that represents that the form is submitted and the items pass the validators
type FormSubmittedMsg struct {
	Data map[string]string
}

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

	// Default values
	Default map[string]string
}

// Initializes a new form item
func New() Form {
	return Form{}
}

// Adds a new input to the form
func (form *Form) AddInput(id, title, desc string, input textinput.Model, validators []Validator) {
	form.Items = append(form.Items, FormItemWrapped{
		ID: id,
		FormItem: &FormItemInputBox{
			Title:       title,
			Description: desc,
			Input:       input,
		},
		Enabled:    true,
		Validators: validators,
	})
}

// Adds a new input to the form
func (form *Form) AddOption(id, title, desc string, options []string) {
	form.Items = append(form.Items, FormItemWrapped{
		ID: id,
		FormItem: &FormItemOptionBox{
			Title:       title,
			Description: desc,
			Values:      options,
		},
		Enabled:    true,
		Validators: []Validator{},
	})
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
func (form *Form) Update(msg tea.Msg, vault *tlockvault.Vault) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
	match_key:
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
		case "enter":
			data := make(map[string]string)

			// Validate them all!
			for _, item := range form.Items {
				// Get the value
				value := strings.TrimSpace(item.FormItem.Value())

				// If it is empty, lets replace with default value
				if defaultValue, ok := form.Default[item.ID]; ok && value == "" {
					value = defaultValue
				}

				// Remove current error
				item.FormItem.SetError(nil)

				// Run validators
				for _, validator := range item.Validators {
					// Validate
					if err := validator(vault, value); err != nil {
						// Set erro
						item.FormItem.SetError(&err)

						// There is issue with the form, break from the switch
						break match_key
					}
				}

				// Set the item
				data[item.ID] = value
			}

			// Return
			cmds = append(cmds, func() tea.Msg { return FormSubmittedMsg{Data: data} })
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

// Runs the post init hook
func (form *Form) PostInit() {
	// Let us focus the first item
	form.Items[0].FormItem.Focus()
}

// Disables a form item
func (form *Form) Disable(id string) {
	// Find the index
	index := slices.IndexFunc(form.Items, func(item FormItemWrapped) bool { return item.ID == id })

	// Disable it
	form.Items[index].Enabled = false
}

// Enables a form item
func (form *Form) Enable(id string) {
	// Find the index
	index := slices.IndexFunc(form.Items, func(item FormItemWrapped) bool { return item.ID == id })

	// Enable it
	form.Items[index].Enabled = true
}
