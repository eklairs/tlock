package folders

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/messages"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

var editFolderAscii = `
█▀▀ █▀▄ █ ▀█▀
██▄ █▄▀ █  █`

// Edit folder key map
type editFolderKeyMap struct {
	GoBack key.Binding
	Enter  key.Binding
}

// ShortHelp()
func (k editFolderKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.GoBack, k.Enter}
}

// FullHelp()
func (k editFolderKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.GoBack},
		{k.Enter},
	}
}

// Keys
var editFolderKeys = editFolderKeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "edit folder"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

// Edit folder screen
type EditFolderScreen struct {
	// Folder name input
	name textinput.Model

	// Foldet to be edited
	folder tlockvault.Folder

	// Error message
	errorMessage *string

	// Vault
	vault *tlockvault.Vault
}

// Initialize edit folder screen
func InitializeEditFolderScreen(folder tlockvault.Folder, vault *tlockvault.Vault) EditFolderScreen {
	// Initialize input box
	name := components.InitializeInputBox("Your folder name goes here...")
	name.SetValue(folder.Name)
	name.Focus()

	// Return
	return EditFolderScreen{
		name:   name,
		folder: folder,
		vault:  vault,
	}
}

// Init
func (screen EditFolderScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen EditFolderScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, editFolderKeys.GoBack):
			manager.PopScreen()

		case key.Matches(msgType, editFolderKeys.Enter):
			// If it is empty, then show error
			if screen.name.Value() == "" {
				screen.errorMessage = &ERROR_EMPTY_FOLDER_NAME
				break
			}

			// Update the folder
			screen.vault.RenameFolder(screen.folder.ID, screen.name.Value())

			// Request folders refresh
			cmds = append(cmds, func() tea.Msg { return tlockmessages.RefreshFoldersMsg{} })

			// Pop
			manager.PopScreen()
		}
	}

	// Update input box
	screen.name, _ = screen.name.Update(msg)

	// Return
	return screen, tea.Batch(cmds...)
}

// View
func (screen EditFolderScreen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(editFolderAscii), "",
		tlockstyles.Styles.SubText.Render("Rename the folder to a new name"), "",
		components.InputGroup("Name", "Choose the new name for your folder", screen.errorMessage, screen.name),
		tlockstyles.Help.View(editFolderKeys),
	)
}
