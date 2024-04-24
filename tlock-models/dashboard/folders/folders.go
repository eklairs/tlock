package folders

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

type folderListItem struct {
	// Name of the folder
	Name string

	// Number of tokens in the folder
	TokensCount int
}

func (item folderListItem) FilterValue() string {
	return item.Name
}

// Folders list delegate
type folderListDelegate struct{}

// Height
func (d folderListDelegate) Height() int {
	return 4
}

// Spacing
func (d folderListDelegate) Spacing() int {
	return 0
}

// Update
func (d folderListDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

// Render
func (d folderListDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item := listItem.(folderListItem)

	// Decide renderer function
	render_fn := components.InactiveFolderListItem

	if index == m.Index() {
		render_fn = components.ActiveFolderListItem
	}

	// Print
	fmt.Fprint(w, render_fn(item.Name, item.TokensCount))
}

// Folders
type Folders struct {
	// Vault
	vault *tlockvault.TLockVault

	// List view
	listview list.Model
}

// Initializes a new instance of folders
func InitializeFolders(vault *tlockvault.TLockVault) Folders {
	return Folders{
		vault: vault,
	}
}

// Handles update messages
func (folders *Folders) Update(msg tea.Msg, manager *modelmanager.ModelManager) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "A":
			manager.PushScreen(InitializeAddFolderScreen())
		case "E":
			// Get focused folder item
			focused_folder := folders.listview.Items()[folders.listview.Index()].(folderListItem)

			// Move to edit screen
			manager.PushScreen(InitializeEditFolderScreen(focused_folder.Name))
		}
	}

	return tea.Batch(cmds...)
}

// View
func (folders Folders) View() string {
	return folders.listview.View()
}
