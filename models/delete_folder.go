package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tlockvault "github.com/eklairs/tlock/tlock-vault"

	. "github.com/eklairs/tlock/internal/modelmanager"
)

var DELETE_FOLDER_SIZE = 60

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


// Delete folder Model
type DeleteFolderModel struct {
    // Styles
    styles Styles

    // Vault
    vault *tlockvault.TLockVault

    // Name of the folder to be deleted
    folder string

    // Help
    help help.Model
}

// Initialize root model
func InitializeDeleteFolderModel(vault *tlockvault.TLockVault, folder string) DeleteFolderModel {
    return DeleteFolderModel {
        styles: InitializeStyles(60),
        vault: vault,
        folder: folder,
        help: help.New(),
    }
}

// Init
func (m DeleteFolderModel) Init() tea.Cmd {
    return nil
}

// Update
func (m DeleteFolderModel) Update(msg tea.Msg, manager *ModelManager) (Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch msgType.String() {
        case "esc":
            manager.PopScreen()
        case "enter":
            m.vault.DeleteFolder(m.folder)

            manager.PopScreen()
        }
    }

	return m, nil
}


// View
func (m DeleteFolderModel) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Left,
        m.styles.dimmedCenter.Render("Are you sure that you want to delete"),
        m.styles.center.Render(lipgloss.JoinHorizontal(
            lipgloss.Center,
            m.styles.title.Copy().UnsetWidth().Render(fmt.Sprintf(" Folder %s", m.folder)),
            m.styles.dimmed.Copy().UnsetWidth().Render("? This action is "),
            m.styles.error.Copy().UnsetWidth().Render("irreversible"),
        )), "",
        m.styles.center.Render(m.help.View(deleteFolderKeys)),
    )
}

