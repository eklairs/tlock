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

var ADD_FOLDER_SIZE = 65

type addFolderKeyMap struct {
    GoBack key.Binding
    Enter key.Binding
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
var addFolderKeys =  addFolderKeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "create folder"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

// Root Model
type AddFolderModel struct {
    // Styles
    styles Styles

    // Folder name input
    name textinput.Model

    // Vault
    vault *tlockvault.TLockVault

    // Help
    help help.Model
}

// Initialize add folder model
func InitializeAddFolderModel(vault *tlockvault.TLockVault) AddFolderModel {
    name := InitializeInputBox("Your folder name goes here...")
    name.Focus()

    return AddFolderModel {
        name: name,
        vault: vault,
        help: help.New(),
        styles: InitializeStyles(ADD_FOLDER_SIZE),
    }
}

// Init
func (m AddFolderModel) Init() tea.Cmd {
    return nil
}

// Update
func (m AddFolderModel) Update(msg tea.Msg, manager *ModelManager) (Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch msgType.String() {
        case "enter":
            m.vault.AddFolder(m.name.Value())

            manager.PopScreen()
        }
    }

    m.name, _ = m.name.Update(msg)

	return m, nil
}

// View
func (m AddFolderModel) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Left,
        m.styles.title.Render("Name"), // Folder name header
        m.styles.dimmed.Render("Choose a name for your folder, like Socials!"), "", // Folder name description
        m.styles.input.Render(m.name.View()), "",
        m.styles.center.Render(m.help.View(addFolderKeys)),
    )
}
