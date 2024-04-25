package dashboard

import (
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-models/dashboard/folders"
	"github.com/eklairs/tlock/tlock-models/dashboard/tokens"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"golang.org/x/term"
)

var EmptyAsciiArt = `
\    /\
 )  ( ')
(  /  )
 \(__)|
`

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
	vault *tlockvault.TLockVault

	// Folders
	folders folders.Folders

	// Tokens
	tokens tokens.Tokens

	// Help
	help help.Model
}

// Initializes a new instance of dashboard screen
func InitializeDashboardScreen(vault tlockvault.TLockVault, context context.Context) DashboardScreen {
	return DashboardScreen{
		vault:   &vault,
		help:    components.BuildHelp(),
		tokens:  tokens.InitializeTokens(&vault, context),
		folders: folders.InitializeFolders(&vault),
	}
}

// Init
func (screen DashboardScreen) Init() tea.Cmd {
	var cmd tea.Cmd

	if len(screen.vault.Data.Folders) != 0 {
		cmd = func() tea.Msg {
			return folders.FolderChangedMsg{
				Folder: screen.vault.Data.Folders[0].Name,
			}
		}
	}

	return tea.Batch(cmd, screen.tokens.Init())
}

// Update
func (screen DashboardScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	return screen, tea.Batch(screen.folders.Update(msg, manager), screen.tokens.Update(msg, manager))
}

// View
func (screen DashboardScreen) View() string {
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))

	if len(screen.vault.Data.Folders) == 0 {
		style := lipgloss.NewStyle().
			Height(height).
			Align(lipgloss.Center, lipgloss.Center)

		ui := lipgloss.JoinVertical(
			lipgloss.Center,
			tlockstyles.Styles.Title.Render(EmptyAsciiArt),
			tlockstyles.Styles.SubText.Render("So empty! How about adding a new folder?"), "",
			screen.help.View(dashboardKeys),
		)

		return style.Render(ui)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		screen.folders.View(), "  ",
		screen.tokens.View(),
	)
}
