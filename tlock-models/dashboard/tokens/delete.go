package tokens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/messages"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

var deleteTokenAsciiArt = `
█▀▄ █▀▀ █   █▀▀ ▀█▀ █▀▀
█▄▀ ██▄ █▄▄ ██▄  █  ██▄`

// Delete token model key bindings
type deleteTokenKeyMap struct {
	Delete key.Binding
	GoBack key.Binding
}

// ShortHelp()
func (k deleteTokenKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.GoBack, k.Delete}
}

// FullHelp()
func (k deleteTokenKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.GoBack},
		{k.Delete},
	}
}

// Keys
var deleteTokenKeys = deleteTokenKeyMap{
	Delete: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "delete"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

// Delete token screen
type DeleteTokenScreen struct {
	// Vault
	vault *tlockvault.Vault

	// Token to delete
	token tlockvault.Token

	// Folder in which the token is
	folder tlockvault.Folder
}

// Initialize root model
func InitializeDeleteTokenScreen(vault *tlockvault.Vault, folder tlockvault.Folder, token tlockvault.Token) DeleteTokenScreen {
	// Return
	return DeleteTokenScreen{
		folder: folder,
		token:  token,
		vault:  vault,
	}
}

// Init
func (screen DeleteTokenScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen DeleteTokenScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, deleteTokenKeys.GoBack):
			manager.PopScreen()
		case key.Matches(msgType, deleteTokenKeys.Delete):
			// Delete
			screen.vault.DeleteToken(screen.folder.Name, screen.token)

			accountName := screen.token.Account

			if accountName == "" {
				accountName = "<no account name>"
			}

			// Require refresh of folders and tokens list
			cmds = append(
				cmds,
				func() tea.Msg { return tlockmessages.RefreshFoldersMsg{} },
				func() tea.Msg { return tlockmessages.RefreshTokensMsg{} },
				func() tea.Msg {
					return components.StatusBarMsg{Message: fmt.Sprintf("Successfully deleted the token (%s)", accountName)}
				},
			)

			// Pop
			manager.PopScreen()
		}
	}

	return screen, tea.Batch(cmds...)
}

// View
func (screen DeleteTokenScreen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(deleteTokenAsciiArt), "",
		tlockstyles.Styles.SubText.Render("Permanently delete token"), "",
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			tlockstyles.Styles.SubText.Render("Are you sure you want to "),
			tlockstyles.Styles.Error.Render("× DELETE "),
			tlockstyles.Styles.Title.Render(screen.token.Account),
			tlockstyles.Styles.SubText.Render(" token ?"),
		), "",
		tlockstyles.Help.View(deleteTokenKeys),
	)
}
