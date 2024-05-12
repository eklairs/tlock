package tokens

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eklairs/tlock/tlock-internal/form"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

func value(input textinput.Model, value string) textinput.Model {
	input.SetValue(value)

	return input
}

// Edit token key map
type editTokenKeyMap struct {
	GoBack key.Binding
	Enter  key.Binding
	Tab    key.Binding
	Arrow  key.Binding
}

// ShortHelp()
func (k editTokenKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Arrow, k.Enter, k.GoBack}
}

// FullHelp()
func (k editTokenKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

// Keys
var editTokenKeys = editTokenKeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "edit token"),
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

// Edit token ascii art
var editTokenAscii = `
█▀▀ █▀▄ █ ▀█▀
██▄ █▄▀ █  █`

// Edit token desc
var editTokenDesc = "Edit your token [secret is required, rest are optional]"

// Edit token screen
type EditTokenScreen struct {
	// Form
	form form.Form

	// Vault
	vault *tlockvault.Vault

	// Folder
	folder tlockvault.Folder

	// Token to edit
	token tlockvault.Token

	// Viewport
	viewport viewport.Model

	// Viewport content
	content string
}

// Initializes a new screen of EditTokenScreen
func InitializeEditTokenScreen(folder tlockvault.Folder, token tlockvault.Token, vault *tlockvault.Vault) EditTokenScreen {
    // Form
    form := BuildForm()

    // Return
	return EditTokenScreen{
		form:     form,
		token:    token,
		vault:    vault,
		folder:   folder,
		viewport: IntoViewport(editTokenAscii, editTokenDesc, form),
	}
}

// Init
func (screen EditTokenScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen EditTokenScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
    // Commands
    cmds := make([]tea.Cmd, 0)

    // Let the form handle
    screen.form.Update(msg, screen.vault)

    // Return
	return screen, tea.Batch(cmds...)
}

// View
func (screen EditTokenScreen) View() string {
	return screen.viewport.View()
}

