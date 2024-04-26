package dashboard

import (
	tea "github.com/charmbracelet/bubbletea"
	tlockvault "github.com/eklairs/tlock/tlock-vault"

	"github.com/charmbracelet/bubbles/key"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
)

// Dashboard key map
type dashboardKeyMap struct {
	Help        key.Binding
	Add         key.Binding
	ChangeTheme key.Binding
}

// ShortHelp()
func (k dashboardKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Add, k.ChangeTheme}
}

// FullHelp()
func (k dashboardKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help},
		{k.Add},
		{k.ChangeTheme},
	}
}

// Keys
var dashboardKeys = dashboardKeyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help menu"),
	),
	Add: key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "add folder"),
	),
	ChangeTheme: key.NewBinding(
		key.WithKeys("ctrl+t"),
		key.WithHelp("Ctrl + T", "change theme"),
	),
}

// Dashboard screen
type DashboardScreen struct {
	// Vault
	vault *tlockvault.Vault

    // Context
    context context.Context
}

// Initializes a new instance of dashboard screen
func InitializeDashboardScreen(vault tlockvault.Vault, context context.Context) DashboardScreen {
	return DashboardScreen{
		vault:   &vault,
        context: context,
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
    return ""
}
