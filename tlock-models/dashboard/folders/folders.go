package folders

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"golang.org/x/term"
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

// Builds the listview for the given list of folders
func buildListViewForFolders(vault *tlockvault.TLockVault) list.Model {
	items := make([]list.Item, len(vault.Data.Folders))

	for index, folder := range vault.Data.Folders {
		items[index] = folderListItem{
			Name:        folder.Name,
			TokensCount: len(folder.Uris),
		}
	}

	return components.ListViewSimple(items, folderListDelegate{}, 65, 10)
}

// Initializes a new instance of folders
func InitializeFolders(vault *tlockvault.TLockVault) Folders {
	return Folders{
		vault:    vault,
		listview: buildListViewForFolders(vault),
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

	folders.listview, _ = folders.listview.Update(msg)

	return tea.Batch(cmds...)
}

// View
func (folders Folders) View() string {
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))

	// Build UI
	ui := lipgloss.JoinVertical(
		lipgloss.Left,
		tlockstyles.Styles.TitleBar.Render("FOLDERS"), "",
		folders.listview.View(),
	)

	// Style
	style := lipgloss.NewStyle().Height(height)

	// Render
	return style.Render(ui)
}
