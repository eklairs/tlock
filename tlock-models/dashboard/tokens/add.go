package tokens

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

var addTokenAscii = `
▄▀█ █▀▄ █▀▄
█▀█ █▄▀ █▄▀`

// Add token key map
type addTokenKeyMap struct {
	GoBack key.Binding
	Enter  key.Binding
	Tab    key.Binding
	Arrow  key.Binding
}

// ShortHelp()
func (k addTokenKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Arrow, k.Enter, k.GoBack}
}

// FullHelp()
func (k addTokenKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

// Keys
var addTokenKeys = addTokenKeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "create token"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab", "shift+tab"),
		key.WithHelp("tab/shift+tab", "switch input"),
	),
	Arrow: key.NewBinding(
		key.WithKeys("right", "left"),
		key.WithHelp("→/←/right/left", "change option"),
	),
}

// Indexes
const (
	ACCOUNT_NAME_INDEX    = 0
	ISSUER_NAME_INDEX     = 1
	SECRET_INDEX          = 2
	TYPE_INDEX            = 3
	HASH_FN_INDEX         = 4
	PERIOD_INDEX          = 5
	INITIAL_COUNTER_INDEX = 6
	DIGITS_INDEX          = 7
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
				renderables[index] = tlockstyles.Styles.SubAltBgItem.Render(option)
			}
		} else {
			renderables[index] = tlockstyles.Styles.SubTextBgItem.Render(option)
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

type AddTokenScreen struct {
	// Items
	items []FormItem

	// Focused index
	focused_index int
}

// Initializes a new instance of the add token screen
func InitializeAddTokenScreen() AddTokenScreen {
	return AddTokenScreen{
		items: []FormItem{
			FormItemInputBox{
				title:       "Account name",
				description: "Name of the account, like John Doe",
				input:       components.InitializeInputBox("Your account name goes here..."),
			},
			FormItemInputBox{
				title:       "Issuer",
				description: "Name of the issuer, like GitHub",
				input:       components.InitializeInputBox("Your issuer name goes here..."),
			},
			FormItemInputBox{
				title:       "Secret",
				description: "Secret provided by the issuer",
				input:       components.InitializeInputBox("The secret goes here..."),
			},
			FormItemOptionBox{
				title:         "Type",
				description:   "Time or counter based token",
				values:        []string{"TOTP", "HOTP"},
				selectedIndex: 0,
			},
			FormItemOptionBox{
				title:         "Hash function",
				description:   "Hash function for the token",
				values:        []string{"SHA1", "SHA256", "SHA512", "MD5"},
				selectedIndex: 1,
			},
			FormItemInputBox{
				title:       "Period",
				description: "Time to refresh the token",
				input:       components.InitializeInputBoxCustomWidth("Time in seconds...", 24),
			},
			FormItemInputBox{
				title:       "Initial counter",
				description: "Initial counter for HOTP token",
				input:       components.InitializeInputBoxCustomWidth("Initial counter goes here...", 24),
			},
			FormItemInputBox{
				title:       "Digits",
				description: "Number of digits",
				input:       components.InitializeInputBoxCustomWidth("Number of digits goes here...", 24),
			},
		},
	}
}

// Init
func (screen AddTokenScreen) Init() tea.Cmd {
	screen.items[ACCOUNT_NAME_INDEX] = screen.items[ACCOUNT_NAME_INDEX].Focus()

	return nil
}

// Switches focus from one index to another
func (screen *AddTokenScreen) switchFocus(old, new int) {
	screen.items[old] = screen.items[old].Unfocus()
	screen.items[new] = screen.items[new].Focus()

	screen.focused_index = new
}

// Update
func (screen AddTokenScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	// Update the current focused item
	screen.items[screen.focused_index], _ = screen.items[screen.focused_index].Update(msg)

	// Handle switch
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		typeInfo, _ := screen.items[TYPE_INDEX].(FormItemOptionBox)

		switch msgType.String() {
		case "tab":
			if screen.focused_index != len(screen.items)-1 {
				next := screen.focused_index + 1

				// If the focused index is previous to the hash function, then decide wether to jump to period or initial counter based on type of token
				if screen.focused_index == HASH_FN_INDEX {
					if typeInfo.SelectedValue() == "TOTP" {
						next = PERIOD_INDEX
					} else {
						next = INITIAL_COUNTER_INDEX
					}
				}

				screen.switchFocus(screen.focused_index, next)
			}
		case "shift+tab":
			if screen.focused_index != 0 {
				next := screen.focused_index - 1

				if screen.focused_index == PERIOD_INDEX || screen.focused_index == INITIAL_COUNTER_INDEX {
					next = HASH_FN_INDEX
				}

				screen.switchFocus(screen.focused_index, next)
			}
		}
	}

	return screen, nil
}

// View
func (screen AddTokenScreen) View() string {
	items := []string{
		tlockstyles.Styles.Title.Render(addTokenAscii), "",
		tlockstyles.Styles.SubText.Render("Add a new token"), "",
		screen.items[ACCOUNT_NAME_INDEX].View(),
		screen.items[ISSUER_NAME_INDEX].View(),
		screen.items[SECRET_INDEX].View(),
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			screen.items[TYPE_INDEX].View(), "   ",
			screen.items[HASH_FN_INDEX].View(),
		), "",
	}

	// Item for type
	typeInfo, _ := screen.items[TYPE_INDEX].(FormItemOptionBox)

	if typeInfo.SelectedValue() == "TOTP" {
		ui := lipgloss.JoinHorizontal(
			lipgloss.Left,
			screen.items[PERIOD_INDEX].View(), "   ",
			screen.items[DIGITS_INDEX].View(),
		)

		items = append(items, ui)
	} else {
		ui := lipgloss.JoinHorizontal(
			lipgloss.Left,
			screen.items[INITIAL_COUNTER_INDEX].View(), "   ",
			screen.items[DIGITS_INDEX].View(),
		)

		items = append(items, ui)
	}

	// Add help
	items = append(items, tlockstyles.Help.View(addTokenKeys))

	return lipgloss.JoinVertical(
		lipgloss.Center,
		items...,
	)
}
