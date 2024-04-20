package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/internal/boundedinteger"
	"github.com/eklairs/tlock/internal/modelmanager"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"golang.org/x/term"
)

var FOLDERS_WIDTH = 30

// Folders
type Folders struct {
    // Vault
    vault tlockvault.TLockVault

    // Focused index
    focused_index boundedinteger.BoundedInteger

    // Styles
    styles Styles
}

// Initializes a new instance of folders
func InitializeFolders(vault tlockvault.TLockVault) Folders {
    return Folders {
        vault: vault,
        styles: InitializeStyles(FOLDERS_WIDTH),
        focused_index: boundedinteger.New(0, len(vault.Data.Folders)),
    }
}

// Handles update messages
func (folders *Folders) Update(msg tea.Msg, manager *modelmanager.ModelManager) tea.Cmd {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch msgType.String() {
        case "J":
            folders.focused_index.Increase()
        case "K":
            folders.focused_index.Decrease()
        case "X":
            manager.PushScreen(InitializeDeleteFolderModel(&folders.vault, folders.vault.Data.Folders[folders.focused_index.Value].Name))
        case "shift+down":
            folders.vault.MoveFolderDown(folders.vault.Data.Folders[folders.focused_index.Value].Name)
            folders.focused_index.Increase()
        case "shift+up":
            folders.vault.MoveFolderUp(folders.vault.Data.Folders[folders.focused_index.Value].Name)
            folders.focused_index.Decrease()
        }
    }
    return nil
}

// View
func (folders Folders) View() string {
    _, height, _ := term.GetSize(0)

    style := lipgloss.NewStyle().
        Width(FOLDERS_WIDTH + 1).
        Height(height)

    items := make([]string, len(folders.vault.Data.Folders))

    for index, folder := range folders.vault.Data.Folders {
        render_fn := folders.styles.folderInactive.Render

        if index == folders.focused_index.Value {
            render_fn = folders.styles.folderActive.Render
        }

        ui := lipgloss.JoinVertical(
            lipgloss.Left,
            folders.styles.title.Render(folder.Name),
            folders.styles.dimmed.Render(fmt.Sprintf("%d tokens", len(folder.Uris))),
        )

        items[index] = render_fn(ui)
    }

    return style.Render(lipgloss.JoinVertical(lipgloss.Left, items...))
}

