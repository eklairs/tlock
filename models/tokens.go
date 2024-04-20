package models

import (
	tea "github.com/charmbracelet/bubbletea"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

// Tokens
type Tokens struct {
    vault tlockvault.TLockVault
}

// Initializes a new instance of folders
func InitializeTokens(vault tlockvault.TLockVault) Tokens {
    return Tokens {
        vault: vault,
    }
}

// Handles update messages
func (folders *Tokens) Update(msg tea.Msg) tea.Cmd {
    return nil
}

// View
func (folders Tokens) View() string {
    return "Tokens"
}

