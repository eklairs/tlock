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
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/context"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/messages"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-internal/utils"
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

	// Context
	context *context.Context
}

// Returns the folders in the form of list item
func buildFolderListItems(vault *tlockvault.Vault) []list.Item {
	// Mapper function
	mapper := func(folder tlockvault.Folder) list.Item {
		return folderListItem(folder)
	}

	// Map folders
	return utils.Map(vault.Folders, mapper)
}

// Builds the listview for the given list of folders
func buildListViewForFolders(vault *tlockvault.Vault, context *context.Context) list.Model {
	// Get size
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))

	// Build listview
	// Height will be auto handled by the view function
	listview := components.ListViewSimple(buildFolderListItems(vault), folderListDelegate{}, foldersWidth(width), 0) // -4 is for the title

	// Use custom keys
	listview.KeyMap.CursorUp = context.Config.Folder.Previous.Binding
	listview.KeyMap.CursorDown = context.Config.Folder.Next.Binding

	// Return listview
	return listview
}

// Initializes a new instance of folders
func InitializeFolders(vault *tlockvault.Vault, context *context.Context) Folders {
	lastFocused := -1

	if len(vault.Folders) != 0 {
		lastFocused = 0
	}

	return Folders{
		vault:       vault,
		context:     context,
		listview:    buildListViewForFolders(vault, context),
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
		switch {
		// Add new folder
		case key.Matches(msgType, folders.context.Config.Folder.Add.Binding):
			cmds = append(cmds, manager.PushScreen(InitializeAddFolderScreen(folders.vault)))

		// Edit focused token
		case key.Matches(msgType, folders.context.Config.Folder.Edit.Binding):
			if focused := folders.Focused(); focused != nil {
				cmds = append(cmds, manager.PushScreen(InitializeEditFolderScreen(*focused, folders.vault)))
			}

		// Delete focused token
		case key.Matches(msgType, folders.context.Config.Folder.Delete.Binding):
			if focused := folders.Focused(); focused != nil {
				cmds = append(cmds, manager.PushScreen(InitializeDeleteFolderScreen(*focused, folders.vault)))
			}

		case key.Matches(msgType, folders.context.Config.Folder.Next.Binding):
			cmds = append(cmds, func() tea.Msg { return tlockmessages.RequestFolderChanged{} })

		case key.Matches(msgType, folders.context.Config.Folder.Previous.Binding):
			cmds = append(cmds, func() tea.Msg { return tlockmessages.RequestFolderChanged{} })

		// Move folder down
		case key.Matches(msgType, folders.context.Config.Folder.MoveUp.Binding):
			if focused := folders.Focused(); focused != nil {
				if folders.vault.MoveFolderUp(focused.ID) {
					// Move cursor down
					folders.listview.CursorUp()

					// Refresh
					cmds = append(
						cmds,
						func() tea.Msg { return tlockmessages.RefreshFoldersMsg{} },
						func() tea.Msg {
							return components.StatusBarMsg{Message: fmt.Sprintf("Successfully moved %s folder up", focused.Name)}
						},
					)
				}
			}

		// Move folder down
		case key.Matches(msgType, folders.context.Config.Folder.MoveDown.Binding):
			if focused := folders.Focused(); focused != nil {
				if folders.vault.MoveFolderDown(focused.ID) {
					// Move cursor down
					folders.listview.CursorDown()

					// Refresh
					cmds = append(
						cmds,
						func() tea.Msg { return tlockmessages.RefreshFoldersMsg{} },
						func() tea.Msg {
							return components.StatusBarMsg{Message: fmt.Sprintf("Successfully moved %s folder up", focused.Name)}
						},
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
		folders.listview.SetHeight(msgType.Height - 6)

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

	// Set the height of the listview
	folders.listview.SetHeight(height - 5)

	// Build UI
	ui := lipgloss.JoinVertical(
		lipgloss.Left, "",
		tlockstyles.Styles.AccentBgItem.Render("FOLDERS"), "",
		folders.listview.View(),
	)

	// Style
	// Lets set the height-1 for the bottom bar
	style := lipgloss.NewStyle().Height(height - 2)

	// Render
	return style.Render(ui)
}
