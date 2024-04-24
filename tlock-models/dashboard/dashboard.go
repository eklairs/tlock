package dashboard

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-models/dashboard/folders"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

// Dashboard screen
type DashboardScreen struct {
	// Vault
	vault *tlockvault.TLockVault

	// Folders
	folders folders.Folders
}

// Initializes a new instance of dashboard screen
func InitializeDashboardScreen(vault tlockvault.TLockVault) DashboardScreen {
	return DashboardScreen{
		vault:   &vault,
		folders: folders.InitializeFolders(&vault),
	}
}

// Init
func (screen DashboardScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen DashboardScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	return screen, tea.Batch(screen.folders.Update(msg, manager))
}

// View
func (screen DashboardScreen) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		screen.folders.View(),
	)
}
