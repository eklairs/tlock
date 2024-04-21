package dashboard

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-models/dashboard/folders"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

// Dashboard Model
type DashboardModel struct {
	// Vault
	vault tlockvault.TLockVault

    // Folders
    folders folders.Folders
}

func InitializeDashboardModel(vault tlockvault.TLockVault, context context.Context) DashboardModel {
	return DashboardModel{
		vault: vault,
        folders: folders.InitializeFolders(vault, context),
	}
}

// Init
func (m DashboardModel) Init() tea.Cmd {
	return nil
}

// Update
func (m DashboardModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	return m, tea.Batch(m.folders.Update(msg, manager))
}

// View
func (m DashboardModel) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
        m.folders.View(),
	)
}
