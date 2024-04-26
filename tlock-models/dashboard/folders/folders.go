package folders

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tlockinternal "github.com/eklairs/tlock/tlock-internal"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"golang.org/x/term"
)

type folderListItem tlockvault.Folder

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
	fmt.Fprint(w, render_fn(item.Name, len(item.Tokens)))
}

// Folders
type Folders struct {
	// Vault
	vault *tlockvault.Vault

	// List view
	listview list.Model
}

// Returns the folders in the form of list item
func buildFolderListItems(vault *tlockvault.Vault) []list.Item {
    // Mapper function
    mapper := func(folder tlockvault.Folder) list.Item {
        return folderListItem(folder)
    }

    // Map folders
	return tlockinternal.Map(vault.Folders, mapper)
}

// Builds the listview for the given list of folders
func buildListViewForFolders(vault *tlockvault.Vault) list.Model {
	// Get size
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))

	// Build listview
	listview := components.ListViewSimple(buildFolderListItems(vault), folderListDelegate{}, 65, height-2) // -2 is for the title

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
func InitializeFolders(vault *tlockvault.Vault) Folders {
	return Folders{
		vault:    vault,
		listview: buildListViewForFolders(vault),
	}
}

// Returns the focused folder item
func (folders Folders) Focused() *tlockvault.Folder {
    // If there are no items, return nil
    if len(folders.listview.Items()) == 0 {
        return nil
    }

    // Get the focused item
    focusedItem := tlockvault.Folder(folders.listview.Items()[folders.listview.Index()].(folderListItem))

    // Return
    return &focusedItem
}

// Handles update messages
func (folders *Folders) Update(msg tea.Msg, manager *modelmanager.ModelManager) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch msgType.String() {
        // Add new folder
        case "A":
            cmds = append(cmds, manager.PushScreen(InitializeAddFolderScreen(folders.vault)))

        // Edit focused token
        case "E":
            if focused := folders.Focused(); focused != nil {
                cmds = append(cmds, manager.PushScreen(InitializeEditFolderScreen(*focused, folders.vault)))
            }

        // Delete focused token
        case "D":
            if focused := folders.Focused(); focused != nil {
                cmds = append(cmds, manager.PushScreen(InitializeDeleteFolderScreen(*focused, folders.vault)))
            }

        // Move folder down
        case "ctrl+up":
            if focused := folders.Focused(); focused != nil {
                if folders.vault.MoveFolderUp(focused.ID) {
                    // Refresh
                    cmds = append(cmds, func() tea.Msg { return tlockinternal.RefreshFoldersMsg{} })

                    // Move cursor up
                    folders.listview.CursorUp()
                }
            }

        // Move folder down
        case "ctrl+down":
            if focused := folders.Focused(); focused != nil {
                if folders.vault.MoveFolderDown(focused.ID) {
                    // Refresh
                    cmds = append(cmds, func() tea.Msg { return tlockinternal.RefreshFoldersMsg{} })

                    // Move cursor down
                    folders.listview.CursorDown()
                }
            }
        }

    // Update items on refresh folders message
    case tlockinternal.RefreshFoldersMsg:
        cmds = append(cmds, folders.listview.SetItems(buildFolderListItems(folders.vault)))
    }

    // Update list
    folders.listview, _ = folders.listview.Update(msg)

    // Return
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

