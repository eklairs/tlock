package dashboard

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

// Dashboard Model
type DashboardModel struct {
	// Vault
	vault tlockvault.TLockVault
}

func InitializeDashboardModel(vault tlockvault.TLockVault) DashboardModel {
	return DashboardModel{
		vault: vault,
	}
}

// Init
func (m DashboardModel) Init() tea.Cmd {
	return nil
}

// Update
func (m DashboardModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	return m, nil
}

// View
func (m DashboardModel) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
	)
}
