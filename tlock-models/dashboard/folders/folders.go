package folders

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tlockinternal "github.com/eklairs/tlock/tlock-internal"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/tlock-messages"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"golang.org/x/term"
)

// Returns the folders width based on the given screen's width
func foldersWidth(width int) int {
	return int(math.Floor((1.0 / 5.0) * float64(width)))
}

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
	fmt.Fprint(w, render_fn(m.Width(), item.Name, len(item.Tokens)))
}

// Folders
type Folders struct {
	// Vault
	vault *tlockvault.Vault

	// List view
	listview list.Model

	// Last focused index of the listview
	// Used for calculating if the list item focus has been changed
	lastFocused int
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
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	// Build listview
	listview := components.ListViewSimple(buildFolderListItems(vault), folderListDelegate{}, foldersWidth(width), height-3) // -3 is for the title

	// Use custom keys
	listview.KeyMap.CursorDown = key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "down"),
	)

	listview.KeyMap.CursorUp = key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "up"),
	)

	return listview
}

// Initializes a new instance of folders
func InitializeFolders(vault *tlockvault.Vault) Folders {
	lastFocused := -1

	if len(vault.Folders) != 0 {
		lastFocused = 0
	}

	return Folders{
		vault:       vault,
		listview:    buildListViewForFolders(vault),
		lastFocused: lastFocused,
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
	var cmd tea.Cmd

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

		case "tab", "shift+tab":
			cmds = append(cmds, func() tea.Msg { return tlockmessages.RequestFolderChanged{} })

		// Move folder down
		case "ctrl+up":
			if focused := folders.Focused(); focused != nil {
				if folders.vault.MoveFolderUp(focused.ID) {
					// Move cursor down
					folders.listview.CursorUp()

					// Refresh
					cmds = append(
						cmds,
						func() tea.Msg { return tlockmessages.RefreshFoldersMsg{} },
						func() tea.Msg { return tlockmessages.RequestFolderChanged{} },
					)
				}
			}

		// Move folder down
		case "ctrl+down":
			if focused := folders.Focused(); focused != nil {
				if folders.vault.MoveFolderDown(focused.ID) {
					// Move cursor down
					folders.listview.CursorDown()

					// Refresh
					cmds = append(
						cmds,
						func() tea.Msg { return tlockmessages.RefreshFoldersMsg{} },
						func() tea.Msg { return tlockmessages.RequestFolderChanged{} },
					)
				}
			}
		}

	// Update items on refresh folders message
	case tlockmessages.RefreshFoldersMsg:
		// Add
		cmds = append(cmds, folders.listview.SetItems(buildFolderListItems(folders.vault)))

		// If this is the first element, might as well post request for folder changed
		if len(folders.listview.Items()) == 1 {
			cmds = append(cmds, func() tea.Msg { return tlockmessages.RequestFolderChanged{} })
		}

		// Handle terminal resizes
	case tea.WindowSizeMsg:
		folders.listview.SetWidth(foldersWidth(msgType.Width))
		folders.listview.SetHeight(msgType.Height - 3)

	case tlockmessages.RequestFolderChanged:
		// New focused item
		if focused := folders.Focused(); focused != nil {
			cmds = append(cmds, func() tea.Msg {
				return tlockmessages.FolderChanged{
					Folder: *focused,
				}
			})
		}
	}

	// Update list
	folders.listview, cmd = folders.listview.Update(msg)
	cmds = append(cmds, cmd)

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
