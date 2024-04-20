package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eklairs/tlock/internal/boundedinteger"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

var FOLDERS_WIDTH = 30

// Folders
type Folders struct {
    // Vault
    vault tlockvault.TLockVault

    // Focused index
    focused_index boundedinteger.BoundedInteger
}

// Initializes a new instance of folders
func InitializeFolders(vault tlockvault.TLockVault) Folders {
    return Folders {
        vault: vault,
        focused_index: boundedinteger.New(0, len(vault.Data.Folders)),
    }
}

// Handles update messages
func (folders *Folders) Update(msg tea.Msg) tea.Cmd {
    return nil
}

// View
func (folders Folders) View() string {
    return "Folders"
}

