package tokens

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/form"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

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
		key.WithHelp("→/←", "change option"),
	),
}

var addTokenAscii = `
▄▀█ █▀▄ █▀▄
█▀█ █▄▀ █▄▀`

// Add token screen
type AddTokenScreen struct {
	// Form
	form form.Form
}

// Initializes a new screen of AddTokenScreen
func InitializeAddTokenScreen(vault *tlockvault.Vault) AddTokenScreen {
	items := []form.FormItem{
		form.FormItemInputBox{
			Title:       "Account Name",
			Description: "Name of the account, like John Doe",
			Input:       components.InitializeInputBox("Account name goes here..."),
		},
		form.FormItemInputBox{
			Title:       "Issuer Name",
			Description: "Name of the issuer, like GitHub",
			Input:       components.InitializeInputBox("Issuer name goes here..."),
		},
		form.FormItemInputBox{
			Title:       "Secret",
			Description: "The secret provided by the issuer",
			Input:       components.InitializeInputBox("The secret goes here..."),
		},
		form.FormItemOptionBox{
			Title:         "Type",
			Description:   "Time or counter based token",
			Values:        []string{"TOTP", "HOTP"},
			SelectedIndex: 0,
		},
		form.FormItemOptionBox{
			Title:         "Hash function",
			Description:   "Hash function for the token",
			Values:        []string{"SHA1", "SHA256", "SHA512", "MD5"},
			SelectedIndex: 1,
		},
		form.FormItemInputBox{
			Title:       "Period",
			Description: "Time to refresh the token",
			Input:       components.InitializeInputBoxCustomWidth("Time in seconds...", 24),
		},
		form.FormItemInputBox{
			Title:       "Initial counter",
			Description: "Initial counter for HOTP token",
			Input:       components.InitializeInputBoxCustomWidth("Initial counter...", 24),
		},
		form.FormItemInputBox{
			Title:       "Digits",
			Description: "Number of digits",
			Input:       components.InitializeInputBoxCustomWidth("Number of digits goes here...", 24),
		},
	}

	// Disable the HOTP initially
	form := form.New(items)
	form.Items[6].Enabled = false

	return AddTokenScreen{
		form: form,
	}
}

// Init
func (screen AddTokenScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen AddTokenScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	// Update the form
	screen.form.Update(msg)

	// Enable / Disable items based on the choosen type
	if screen.form.Items[3].FormItem.Value() == "TOTP" {
		// Enable period item
		screen.form.Items[5].Enabled = true

		// Disable initial counter item
		screen.form.Items[6].Enabled = false
	}

	if screen.form.Items[3].FormItem.Value() == "HOTP" {
		// Enable initial counter item
		screen.form.Items[6].Enabled = true

		// Disable period item
		screen.form.Items[5].Enabled = false
	}

	return screen, nil
}

// View
func (screen AddTokenScreen) View() string {
	// Items
	items := []string{
		tlockstyles.Styles.Title.Render(addTokenAscii), "",
		tlockstyles.Styles.SubText.Render("Add a new token"), "",
		screen.form.Items[0].FormItem.View(), // Account name input
		screen.form.Items[1].FormItem.View(), // Issuer name input
		screen.form.Items[2].FormItem.View(), // Secret value input
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			screen.form.Items[3].FormItem.View(), "   ",
			screen.form.Items[4].FormItem.View(),
		), "",
	}

	// Render the input boxes based on choosen type
	inputGroup := lipgloss.JoinHorizontal(
		lipgloss.Left,
		screen.form.Items[5].FormItem.View(), "   ",
		screen.form.Items[7].FormItem.View(),
	)

	if screen.form.Items[3].FormItem.Value() == "HOTP" {
		inputGroup = lipgloss.JoinHorizontal(
			lipgloss.Left,
			screen.form.Items[6].FormItem.View(), "   ",
			screen.form.Items[7].FormItem.View(),
		)
	}

	// Add the help menu
	items = append(items, inputGroup, "", tlockstyles.Help.View(addTokenKeys))

	// Return!
	return lipgloss.JoinVertical(
		lipgloss.Center,
		items...,
	)
}
