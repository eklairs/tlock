package folders

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

var deleteFolderAsciiArt = `
█▀▄ █▀▀ █   █▀▀ ▀█▀ █▀▀
█▄▀ ██▄ █▄▄ ██▄  █  ██▄`

// Create user model key bindings
type deleteFolderKeyMap struct {
	Delete key.Binding
	GoBack key.Binding
}

// ShortHelp()
func (k deleteFolderKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.GoBack, k.Delete}
}

// FullHelp()
func (k deleteFolderKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.GoBack},
		{k.Delete},
	}
}

// Keys
var deleteFolderKeys = deleteFolderKeyMap{
	Delete: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "delete"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

// Delete folder screen
type DeleteFolderScreen struct {
	// Folder to delete
	folder tlockvault.Folder

	// Vault
	vault *tlockvault.Vault
}

// Initialize root model
func InitializeDeleteFolderScreen(folder tlockvault.Folder, vault *tlockvault.Vault) DeleteFolderScreen {
	return DeleteFolderScreen{
		folder: folder,
		vault:  vault,
	}
}

// Init
func (screen DeleteFolderScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen DeleteFolderScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, deleteFolderKeys.GoBack):
			manager.PopScreen()
		case key.Matches(msgType, deleteFolderKeys.Delete):
			// Delete the folder
			screen.vault.DeleteFolder(screen.folder.ID)

			// Request folders refresh
			cmds = append(
				cmds,
				func() tea.Msg { return tlockmessages.RequestFolderChanged{} },
				func() tea.Msg { return tlockmessages.RefreshFoldersMsg{} },
                func() tea.Msg { return components.StatusBarMsg{Message: fmt.Sprintf("Successfully deleted %s folder", screen.folder.Name)} },
			)

			// Pop
			manager.PopScreen()
		}
	}

	return screen, tea.Batch(cmds...)
}

// View
func (screen DeleteFolderScreen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(deleteFolderAsciiArt), "",
		tlockstyles.Styles.SubText.Render("Permanently delete tokens folder"), "",
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			tlockstyles.Styles.SubText.Render("Are you sure you want to "),
			tlockstyles.Styles.Error.Render("× DELETE "),
			tlockstyles.Styles.Title.Render(screen.folder.Name),
			tlockstyles.Styles.SubText.Render("?"),
		), "",
		tlockstyles.Help.View(deleteFolderKeys),
	)
}
