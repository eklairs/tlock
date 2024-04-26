package folders

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/tlock-messages"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

var ERROR_EMPTY_FOLDER_NAME = "Folder name cannot be empty"

var addFolderAscii = `
▄▀█ █▀▄ █▀▄
█▀█ █▄▀ █▄▀`

// Add folder key map
type addFolderKeyMap struct {
	GoBack key.Binding
	Enter  key.Binding
}

// ShortHelp()
func (k addFolderKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.GoBack, k.Enter}
}

// FullHelp()
func (k addFolderKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.GoBack},
		{k.Enter},
	}
}

// Keys
var addFolderKeys = addFolderKeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "create folder"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

// Add folder screen
type AddFolderScreen struct {
	// Folder name input
	name textinput.Model

	// Error
	errorMessage *string

	// Vault
	vault *tlockvault.Vault
}

// Initialize add folder scree
func InitializeAddFolderScreen(vault *tlockvault.Vault) AddFolderScreen {
	// Initialize input box
	name := components.InitializeInputBox("Your folder name goes here...")
	name.Focus()

	// Return
	return AddFolderScreen{
		name:  name,
		vault: vault,
	}
}

// Init
func (screen AddFolderScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen AddFolderScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, addFolderKeys.GoBack):
			manager.PopScreen()

		case key.Matches(msgType, addFolderKeys.Enter):
			// If it is empty, show error message
			if screen.name.Value() == "" {
				screen.errorMessage = &ERROR_EMPTY_FOLDER_NAME

				break
			}

			// Add the folder
			screen.vault.AddFolder(screen.name.Value())

			// Request folders refresh
			cmds = append(cmds, func() tea.Msg { return tlockmessages.RefreshFoldersMsg{} })

			// Pop
			manager.PopScreen()
		}
	}

	// Send the value to input box
	screen.name, _ = screen.name.Update(msg)

	// Return
	return screen, tea.Batch(cmds...)
}

// View
func (screen AddFolderScreen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(addFolderAscii), "",
		tlockstyles.Styles.SubText.Render("Add a new folder"), "",
		components.InputGroup("Name", "Choose a name for your folder, like Socials!", screen.errorMessage, screen.name),
		tlockstyles.Help.View(addFolderKeys),
	)
}
