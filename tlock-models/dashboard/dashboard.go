package dashboard

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

// Dashboard screen
type DashboardScreen struct {
	// Vault
	vault tlockvault.TLockVault
}

// Initializes a new instance of dashboard screen
func InitializeDashboardScreen(vault tlockvault.TLockVault) DashboardScreen {
	return DashboardScreen{
		vault: vault,
	}
}

// Init
func (screen DashboardScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen DashboardScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	return screen, nil
}

// View
func (screen DashboardScreen) View() string {
	return "Dashboard"
}
