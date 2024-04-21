package models

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
	tlockvault "github.com/eklairs/tlock/tlock-vault"

	. "github.com/eklairs/tlock/internal/modelmanager"
)

// Create user model key bindings
type editFolderKeyMap struct {
    Edit key.Binding
    GoBack key.Binding
}

// ShortHelp()
func (k editFolderKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.GoBack, k.Edit}
}

// FullHelp()
func (k editFolderKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
        {k.GoBack},
        {k.Edit},
	}
}

// Keys
var editFolderKeys =  editFolderKeyMap{
	Edit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "edit"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

// Edit folder Model
type EditFolderModel struct {
    // Styles
    styles Styles

    // Vault
    vault *tlockvault.TLockVault

    // Name of the folder to be deleted
    folder string

    // Help
    help help.Model

    // New folder name input
    folderNameInput textinput.Model
}

// Initialize root model
func InitializeEditFolderModel(vault *tlockvault.TLockVault, folder string) EditFolderModel {
    folderNameInput := InitializeInputBox("Choose a new name for the folder...")
    folderNameInput.SetValue(folder)
    folderNameInput.Focus()

    return EditFolderModel {
        styles: InitializeStyles(60),
        vault: vault,
        folder: folder,
        help: help.New(),
        folderNameInput: folderNameInput,
    }
}

// Init
func (m EditFolderModel) Init() tea.Cmd {
    return nil
}

// Update
func (m EditFolderModel) Update(msg tea.Msg, manager *ModelManager) (Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch msgType.String() {
        case "esc":
            manager.PopScreen()
        case "enter":
            // m.vault.EditFolder(m.folder)

            manager.PopScreen()
        }
    }

	return m, nil
}


// View
func (m EditFolderModel) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Left,
        m.styles.dimmedCenter.Render("Choose a new name for your folder"),
        m.styles.title.Render("Name"), // Username header
        m.styles.dimmed.Render("Choose an awesome name for your folder, like Social"), // Username description 
        m.styles.input.Render(m.folderNameInput.View()),
        m.styles.center.Render(m.help.View(editFolderKeys)),
    )
}

