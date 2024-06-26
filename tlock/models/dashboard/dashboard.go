package dashboard

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"github.com/eklairs/tlock/tlock/models/dashboard/folders"
	"github.com/eklairs/tlock/tlock/models/dashboard/tokens"
	tlockstyles "github.com/eklairs/tlock/tlock/styles"
	"golang.org/x/term"

	"github.com/charmbracelet/bubbles/key"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/config"
	"github.com/eklairs/tlock/tlock-internal/context"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/messages"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
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
var dashboardKeys dashboardKeyMap

// Dashboard screen
type DashboardScreen struct {
	// Vault
	vault *tlockvault.Vault

	// Context
	context *context.Context

	// Folders
	folders folders.Folders

	// Tokens
	tokens tokens.Tokens

	// Status bar
	statusbar components.StatusBar
}

// Initializes a new instance of dashboard screen
func InitializeDashboardScreen(username string, vault *tlockvault.Vault, context *context.Context) modelmanager.Screen {
	// Load keybindings for the user
	context.Config = config.LoadUserConfig(username)

	// Initialize dashboard keymap
	dashboardKeys = dashboardKeyMap{
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help menu"),
		),
		Add: key.NewBinding(
			key.WithKeys(context.Config.Folder.Add.Keys()...),
			key.WithHelp(strings.Join(context.Config.Folder.Add.Keys(), "/"), "add folder"),
		),
		ChangeTheme: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+t", "change theme"),
		),
	}

	return DashboardScreen{
		vault:     vault,
		context:   context,
		statusbar: components.NewStatusBar(username),
		folders:   folders.InitializeFolders(vault, context),
		tokens:    tokens.InitializeTokens(vault, context),
	}
}

// Init
func (screen DashboardScreen) Init() tea.Cmd {
	var cmd tea.Cmd

	if len(screen.vault.Folders) != 0 {
		cmd = func() tea.Msg {
			return tlockmessages.FolderChanged{
				Folder: screen.vault.Folders[0],
			}
		}
	}

	return tea.Batch(cmd, tlockmessages.DispatchRefreshTokensValueMsg())
}

// Update
func (screen DashboardScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	var cmd tea.Cmd

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		// Help menu
		case key.Matches(msgType, dashboardKeys.Help):
			cmd = manager.PushScreen(InitializeHelpScreen(screen.context))

		// Themes screen
		case key.Matches(msgType, dashboardKeys.ChangeTheme):
			cmd = manager.PushScreen(InitializeThemesScreen(screen.context))
		}
	}

	return screen, tea.Batch(screen.folders.Update(msg, manager), screen.tokens.Update(msg, manager), cmd, screen.statusbar.Update(msg))
}

// View
func (screen DashboardScreen) View() string {
	// Get the size of the terminal
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))

	if len(screen.vault.Folders) == 0 {
		style := lipgloss.NewStyle().
			Height(height).
			Align(lipgloss.Center, lipgloss.Center)

		ui := lipgloss.JoinVertical(
			lipgloss.Center,
			tlockstyles.Styles.Title.Render(EmptyAsciiArt),
			tlockstyles.Styles.SubText.Render("So empty! How about adding a new folder?"), "",
			tlockstyles.Help.View(dashboardKeys),
		)

		return style.Render(ui)
	}

	ui := lipgloss.JoinHorizontal(
		lipgloss.Left,
		screen.folders.View(), "  ",
		screen.tokens.View(),
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		ui, screen.statusbar.View(),
	)
}
