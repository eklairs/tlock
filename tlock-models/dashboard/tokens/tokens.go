package tokens

import (
	"golang.org/x/term"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/eklairs/tlock/tlock-internal/boundedinteger"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-models/dashboard/folders"

	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

// Tokens
type Tokens struct {
	// Context
	context context.Context

	// Vault
	vault tlockvault.TLockVault

	// Focused index
	focused_index boundedinteger.BoundedInteger

	// Styles
	styles tlockstyles.Styles
}

// Initializes a new instance of folders
func InitializeTokens(vault tlockvault.TLockVault, context context.Context) Tokens {
	// Terminal size
	width, _, _ := term.GetSize(0)

	// Styles
	styles := tlockstyles.InitializeStyle(width-folders.FOLDERS_WIDTH, context.Theme)

	return Tokens{
		vault:   vault,
		styles:  styles,
		context: context,
	}
}

// Handles update messages
func (tokens *Tokens) Update(msg tea.Msg, manager *modelmanager.ModelManager) tea.Cmd {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "j":
			tokens.focused_index.Increase()
		case "k":
			tokens.focused_index.Decrease()
		}
	}

	return nil
}

// View
func (tokens Tokens) View() string {
	// Get term size
	width, height, _ := term.GetSize(0)

	// Full style
	style := lipgloss.NewStyle().
		Width(width - folders.FOLDERS_WIDTH).
		Height(height)

	// List of items
	items := make([]string, 0)

    // Render
	return style.Render(lipgloss.JoinVertical(lipgloss.Center, items...))
}

