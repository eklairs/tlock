package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/internal/modelmanager"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"golang.org/x/term"
    "github.com/pquerna/otp"
)

type dashboardStyles struct {
    title lipgloss.Style
    titleCenter lipgloss.Style
    dimmed lipgloss.Style
    dimmedCenter lipgloss.Style
    input lipgloss.Style
}

// Root Model
type DashboardModel struct {
    styles dashboardStyles
    vault tlockvault.TLockVault
    current_index int
}

// Initialize root model
func InitializeDashboardModel(vault tlockvault.TLockVault) DashboardModel {
    dimmed := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

    return DashboardModel {
        styles: dashboardStyles{
            title: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")),
            titleCenter: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")).Width(30).Align(lipgloss.Center),
            input: lipgloss.NewStyle().Padding(1, 3).Width(30).Background(lipgloss.Color("#1e1e2e")),
            dimmed: dimmed,
            dimmedCenter: dimmed.Width(30).Copy().Align(lipgloss.Center),
        },
        vault: vault,
        current_index: 0,
    }
}

// Init
func (m DashboardModel) Init() tea.Cmd {
    return nil
}

// Update
func (m DashboardModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch msgType.String() {
        case "J":
            m.current_index = (m.current_index + 1) % len(m.vault.Data.Folders)
        case "K":
            if m.current_index == 0 {
                m.current_index = len(m.vault.Data.Folders) - 1
            } else {
                m.current_index -= 1
            }
        }
    }

	return m, nil
}

// View
func (m DashboardModel) View() string {
    width, height, _ := term.GetSize(0)

    style := lipgloss.NewStyle().Width(30).Height(height)

    // Folders
    folders := make([]string, 0)

    for _, folder := range m.vault.Data.Folders {
        ui := lipgloss.JoinVertical(
            lipgloss.Left,
            m.styles.title.Render(folder.Name),
            m.styles.dimmed.Render(fmt.Sprintf("%d tokens", len(folder.Uris))),
        )

        folders = append(folders, lipgloss.NewStyle().MarginTop(1).Render(ui))
    }

    // Tokens
    tokens := make([]string, 0)

    for _, uri := range m.vault.Data.Folders[m.current_index].Uris {
        style := lipgloss.NewStyle().
            Width(width - 30 - 2).
            Padding(1, 3).
            MarginBottom(1).
            Background(lipgloss.Color("#1E1E2E"))

        totp, _ := otp.NewKeyFromURL(uri)

        tokens = append(tokens, style.Render(fmt.Sprintf("%s • %s", totp.AccountName(), totp.Issuer())))
    }

    ui := []string {
        style.Render(lipgloss.JoinVertical(lipgloss.Left, folders...)),
    }

    ui = append(ui, lipgloss.JoinVertical(lipgloss.Left, tokens...))

    return lipgloss.JoinHorizontal(
        lipgloss.Left,
        ui...
    )
}

