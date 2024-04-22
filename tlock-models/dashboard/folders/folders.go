package folders

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/boundedinteger"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"golang.org/x/term"

	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

var FOLDERS_WIDTH = 40

func notifyFolderChanged(folder string) tea.Cmd {
	return func() tea.Msg {
		return FolderChangedMsg{
			Folder: folder,
		}
	}
}

type FolderChangedMsg struct {
	Folder string
}

// Folders
type Folders struct {
	// Context
	context context.Context

	// Vault
	vault *tlockvault.TLockVault

	// Focused index
	focused_index boundedinteger.BoundedInteger

	// Styles
	styles tlockstyles.Styles
}

// Initializes a new instance of folders
func InitializeFolders(vault *tlockvault.TLockVault, context context.Context) Folders {
	styles := tlockstyles.InitializeStyle(FOLDERS_WIDTH, context.Theme)

	return Folders{
		vault:         vault,
		styles:        styles,
		context:       context,
		focused_index: boundedinteger.New(0, len(vault.Data.Folders)),
	}
}

// Handles update messages
func (folders *Folders) Update(msg tea.Msg, manager *modelmanager.ModelManager) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "J":
			folders.focused_index.Increase()

			cmds = append(cmds, notifyFolderChanged(folders.vault.Data.Folders[folders.focused_index.Value].Name))
		case "K":
			folders.focused_index.Decrease()

			cmds = append(cmds, notifyFolderChanged(folders.vault.Data.Folders[folders.focused_index.Value].Name))
		case "A":
			manager.PushScreen(InitializeAddFolderModel(folders.context))
		case "E":
			manager.PushScreen(InitializeEditFolderModel(folders.vault.Data.Folders[folders.focused_index.Value].Name, folders.context))
		case "X":
			manager.PushScreen(InitializeDeleteFolderModel(folders.vault.Data.Folders[folders.focused_index.Value].Name, folders.context))
		}

	case AddNewFolderMsg:
		// Add folder
		folders.vault.AddFolder(msgType.FolderName)

		// Update focused index bounds
		folders.focused_index = boundedinteger.New(len(folders.vault.Data.Folders)-1, len(folders.vault.Data.Folders))

		// Switch
		cmds = append(cmds, notifyFolderChanged(msgType.FolderName))

	case EditNewFolderMsg:
		folders.vault.RenameFolder(msgType.OldName, msgType.NewName)

		cmds = append(cmds, notifyFolderChanged(msgType.NewName))

	case DeleteFolderMsg:
		folders.vault.DeleteFolder(msgType.FolderName)

		folders.focused_index = boundedinteger.New(min(folders.focused_index.Value, len(folders.vault.Data.Folders)-1), len(folders.vault.Data.Folders))

		cmds = append(cmds, notifyFolderChanged(folders.vault.Data.Folders[folders.focused_index.Value].Name))
	}

	return tea.Batch(cmds...)
}

// View
func (folders Folders) View() string {
	// Get term size
	_, height, _ := term.GetSize(0)

	// Full style
	style := lipgloss.NewStyle().
		Width(FOLDERS_WIDTH + 1).
		Height(height - 3)

	// List of items
	items := make([]string, 0)

    // Header
    items = append(items, folders.styles.AccentTitle.Copy().Margin(1).Render("FOLDERS"))

	for index, folder := range folders.vault.Data.Folders {
		render_fn := folders.styles.FolderInactive.Render

		if index == folders.focused_index.Value {
			render_fn = folders.styles.FolderActive.Render
		}

		ui := lipgloss.JoinVertical(
			lipgloss.Left,
			folders.styles.Title.Copy().UnsetWidth().Render(folder.Name),
			folders.styles.Dimmed.Render(fmt.Sprintf("%d tokens", len(folder.Uris))),
		)

        items = append(items, render_fn(ui))
	}

	return style.Render(lipgloss.JoinVertical(lipgloss.Left, items...))
}
