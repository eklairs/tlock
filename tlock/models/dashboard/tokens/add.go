package tokens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/form"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/messages"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
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
	form := BuildForm(map[string]string{})

	// Return
	return AddTokenScreen{
		form:     form,
		vault:    vault,
		folder:   folder,
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

	// Match
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, addTokenKeys.GoBack):
			manager.PopScreen()
		}

	case form.FormSubmittedMsg:
		// Get token
		token := TokenFromFormData(msgType.Data)

		// Make statusbar message
		statusBarMessage := fmt.Sprintf("Successfully added token for %s", token.Account)

		if token.Account == "" {
			statusBarMessage = fmt.Sprintf("Successfully added token (no account name)")
		}

		// Require refresh of folders and tokens list
		cmds = append(
			cmds,
			func() tea.Msg { return tlockmessages.RefreshFoldersMsg{} },
			func() tea.Msg { return tlockmessages.RefreshTokensMsg{} },
			func() tea.Msg { return components.StatusBarMsg{Message: statusBarMessage} },
		)

		// Add
		screen.vault.AddTokenFromToken(screen.folder.Name, token)

		// Break
		manager.PopScreen()
	}

	// Let the form handle its update
	cmds = append(cmds, screen.form.Update(msg, screen.vault))

	// Update the viewport
	DisableBasedOnType(&screen.form)
	UpdateViewport(addTokenAscii, addTokenDesc, &screen.viewport, screen.form)

	// Send the update message to the viewport
	screen.viewport, _ = screen.viewport.Update(msg)

	// Return
	return screen, tea.Batch(cmds...)
}

// View
func (screen AddTokenScreen) View() string {
	return screen.viewport.View()
}
