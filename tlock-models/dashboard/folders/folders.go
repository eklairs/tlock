package folders

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"golang.org/x/term"
)

// Folder changed message
type FolderChangedMsg struct {
	Folder string
}

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

// Returns the folders in the form of list item
func buildFolderListItems(vault *tlockvault.TLockVault) []list.Item {
	items := make([]list.Item, len(vault.Data.Folders))

	for index, folder := range vault.Data.Folders {
		items[index] = folderListItem{
			Name:        folder.Name,
			TokensCount: len(folder.Uris),
		}
	}

	return items
}

// Builds the listview for the given list of folders
func buildListViewForFolders(vault *tlockvault.TLockVault) list.Model {
	// Get size
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))

	// Build listview
	listview := components.ListViewSimple(buildFolderListItems(vault), folderListDelegate{}, 75, height-2) // -2 is for the title

	// Use custom keys
	listview.KeyMap.CursorDown = key.NewBinding(
		key.WithKeys("J"),
		key.WithHelp("J", "down"),
	)

	listview.KeyMap.CursorUp = key.NewBinding(
		key.WithKeys("K"),
		key.WithHelp("K", "up"),
	)

	return listview
}

// Initializes a new instance of folders
func InitializeFolders(vault *tlockvault.TLockVault) Folders {
	return Folders{
		vault:    vault,
		listview: buildListViewForFolders(vault),
	}
}

// Returns the focused folder
func (folders *Folders) Focused() (*folderListItem, error) {
	// If no items found, just skip
	if len(folders.listview.Items()) == 0 {
		return nil, errors.New("No items on the folders listview")
	}

	// Get the focused item
	focused := folders.listview.Items()[folders.listview.Index()].(folderListItem)

	// Return
	return &focused, nil
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
			if focused, err := folders.Focused(); err == nil {
				manager.PushScreen(InitializeEditFolderScreen(focused.Name))
			}
		case "X":
			if focused, err := folders.Focused(); err == nil {
				manager.PushScreen(InitializeDeleteFolderScreen(focused.Name))
			}
		}

	case AddNewFolderMsg:
		// Add folder
		folders.vault.AddFolder(msgType.FolderName)

		// Rebuild list
		cmds = append(cmds, folders.listview.SetItems(buildFolderListItems(folders.vault)))

	case EditFolderMsg:
		// Rename folder
		folders.vault.RenameFolder(msgType.OldName, msgType.NewName)

		// Update list view
		cmds = append(cmds, folders.listview.SetItems(buildFolderListItems(folders.vault)))

	case DeleteFolderMsg:
		// Rename folder
		folders.vault.DeleteFolder(msgType.FolderName)

		// Update list view
		cmds = append(cmds, folders.listview.SetItems(buildFolderListItems(folders.vault)))
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
		tlockstyles.Styles.AccentBgItem.Render("FOLDERS"), "",
		folders.listview.View(),
	)

	// Style
	style := lipgloss.NewStyle().Height(height)

	// Render
	return style.Render(ui)
}
