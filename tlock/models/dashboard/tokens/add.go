package tokens

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eklairs/tlock/tlock-internal/form"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

// Converts the given string type tok token type
func toTokenType(value string) tlockvault.TokenType {
	if value == "HOTP" {
		return tlockvault.TokenTypeHOTP
	}

	return tlockvault.TokenTypeTOTP
}

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

// Add token ascii
var addTokenAscii = `
▄▀█ █▀▄ █▀▄
█▀█ █▄▀ █▄▀`

// Add token description
var addTokenDesc = "Create a new token [only secret is required; rest are optional]"

// Add token screen
type AddTokenScreen struct {
	// Form
	form form.Form

	// Vault
	vault *tlockvault.Vault

	// Folder
	folder tlockvault.Folder

	// Viewport
	viewport viewport.Model

	// Viewport content
	content string
}

// Initializes a new screen of AddTokenScreen
func InitializeAddTokenScreen(folder tlockvault.Folder, vault *tlockvault.Vault) AddTokenScreen {
    // Initialize form
    form := BuildForm()

	// Return
	return AddTokenScreen{
        form: form,
        vault: vault,
        folder: folder,
        viewport: IntoViewport(addTokenAscii, addTokenDesc, form),
    }
}

// Init
func (screen AddTokenScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen AddTokenScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
    // Commands
    cmds := make([]tea.Cmd, 0)

    // Let the form handle its update
    screen.form.Update(msg)

    // Update the viewport
    DisableBasedOnType(&screen.form)
    UpdateViewport(addTokenAscii, addTokenDesc, &screen.viewport, screen.form)

    // Return
	return screen, tea.Batch(cmds...)
}

// View
func (screen AddTokenScreen) View() string {
	return screen.viewport.View()
}

