package models

import (
	"strings"

	"golang.org/x/term"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
	tlockvault "github.com/eklairs/tlock/tlock-vault"

	. "github.com/eklairs/tlock/internal/modelmanager"
)

var EXCLAMATION_MARK = `
┏┓
┃┃
┃┃
┗┛
┏┓
┗┛
`

// Dashboard Model
type DashboardModel struct {
    // Vault
    vault tlockvault.TLockVault

    // Folders
    folders Folders

    // Tokens
    tokens Tokens
}

func InitializeDashboardModel(vault tlockvault.TLockVault) DashboardModel {
    return DashboardModel{
        vault: vault,
        folders: InitializeFolders(vault),
        tokens: InitializeTokens(vault),
    }
}

// Init
func (m DashboardModel) Init() tea.Cmd {
    return nil
}

// Update
func (m DashboardModel) Update(msg tea.Msg, manager *ModelManager) (Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch msgType.String() {
        case "A":
            manager.PushScreen(InitializeAddFolderModel(&m.vault))
        }
    }

    m.folders.Update(msg, manager)
    m.tokens.Update(msg)

    return m, nil
}

// View
func (m DashboardModel) View() string {
    width, height, _ := term.GetSize(0)

    if len(m.vault.Data.Folders) == 0 {
        accent_style := lipgloss.NewStyle().Foreground(COLOR_ACCENT)
        dimmed_style := lipgloss.NewStyle().Foreground(COLOR_DIMMED)

        ui := []string {
            accent_style.Render(EXCLAMATION_MARK),
            "No folders found",
            lipgloss.JoinHorizontal(
                lipgloss.Left,
                dimmed_style.Render("Press "),
                accent_style.Render("A"),
                dimmed_style.Render(" to add a new token"),
            ),
        }

        style := lipgloss.NewStyle().
            Width(width).
            Height(height).
            Align(lipgloss.Center, lipgloss.Center).
            Foreground(COLOR_DIMMED).
            Render(strings.Join(ui, "\n"))

        return style
    }

    return lipgloss.JoinHorizontal(
        lipgloss.Left,
        m.folders.View(),
        m.tokens.View(),
    )
}
