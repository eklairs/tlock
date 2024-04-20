package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/internal/modelmanager"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

var _____ascii = `
 _____   _ _ _   
|   __|_| |_| |_ 
|   __| . | |  _|
|_____|___|_|_|  
`

type deleteTokenStyles struct {
    title lipgloss.Style
    titleCenter lipgloss.Style
    dimmed lipgloss.Style
    dimmedCenter lipgloss.Style
    input lipgloss.Style
}

// Root Model
type DeleteTokenModel struct {
    styles deleteTokenStyles
    vault *tlockvault.TLockVault
    folder int
    token int
}

// Initialize root model
func InitializeDeleteTokenModel(vault *tlockvault.TLockVault, folder, token int) DeleteTokenModel {
    dimmed := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
    styles := deleteTokenStyles{
        title: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")),
        titleCenter: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")).Width(55).Align(lipgloss.Center),
        input: lipgloss.NewStyle().Padding(1, 3).Width(55).Background(lipgloss.Color("#1e1e2e")),
        dimmed: dimmed,
        dimmedCenter: dimmed.Width(55).Copy().Align(lipgloss.Center),
    }

    return DeleteTokenModel {
        styles: styles,
        vault: vault,
        folder: folder,
        token: token,
    }
}

// Init
func (m DeleteTokenModel) Init() tea.Cmd {
    return nil
}

// Update
func (m DeleteTokenModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch msgType.String() {
        case "esc":
            manager.PopScreen()
        case "enter":
            m.vault.DeleteURI(m.folder, m.token)
            manager.PopScreen()
        }
    }

	return m, nil
}

// View
func (m DeleteTokenModel) View() string {
    title := m.styles.title.Copy().UnsetWidth()
    dimmed := m.styles.dimmed.Copy().UnsetWidth()

    return lipgloss.JoinVertical(
        lipgloss.Left,
        m.styles.titleCenter.Render(____ascii),
        m.styles.dimmedCenter.Render("Are you sure you want to delete? This cannot be reversed!"), "",
        lipgloss.JoinHorizontal(
            lipgloss.Center,
            fmt.Sprintf("%s %s %s", dimmed.Render("Press ["), title.Render("Enter"), dimmed.Render("] to continue")),
            "    ",
            fmt.Sprintf("%s %s %s", dimmed.Render("Press ["), title.Render("Esc"), dimmed.Render("] to go back")),
        ),
    )
}

