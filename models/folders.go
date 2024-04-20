package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/internal/boundedinteger"
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
func (folders *Folders) Update(msg tea.Msg) tea.Cmd {
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

        items = append(items, render_fn(ui))
    }

    return style.Render(lipgloss.JoinVertical(lipgloss.Left, items...))
}

