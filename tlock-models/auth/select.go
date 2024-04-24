package auth

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
)

// Select user
type SelectUserScreen struct {
}

// New instance of select user
func InitializeSelectUserScreen() SelectUserScreen {
	return SelectUserScreen{}
}

// Init
func (screen SelectUserScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen SelectUserScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	return screen, nil
}

// View
func (screen SelectUserScreen) View() string {
	return ""
}
