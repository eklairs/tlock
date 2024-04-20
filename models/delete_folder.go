package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/internal/modelmanager"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

var ________ascii = `
 _____   _ _ _   
|   __|_| |_| |_ 
|   __| . | |  _|
|_____|___|_|_|  
`

type deleteFolderStyles struct {
    title lipgloss.Style
    titleCenter lipgloss.Style
    dimmed lipgloss.Style
    dimmedCenter lipgloss.Style
    input lipgloss.Style
}

// Root Model
type DeleteFolderModel struct {
    styles deleteFolderStyles
    vault *tlockvault.TLockVault
    folder int
}

// Initialize root model
func InitializeDeleteFolderModel(vault *tlockvault.TLockVault, folder int) DeleteFolderModel {
    dimmed := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
    styles := deleteFolderStyles{
        title: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")),
        titleCenter: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")).Width(55).Align(lipgloss.Center),
        input: lipgloss.NewStyle().Padding(1, 3).Width(55).Background(lipgloss.Color("#1e1e2e")),
        dimmed: dimmed,
        dimmedCenter: dimmed.Width(55).Copy().Align(lipgloss.Center),
    }

    return DeleteFolderModel {
        styles: styles,
        vault: vault,
        folder: folder,
    }
}

// Init
func (m DeleteFolderModel) Init() tea.Cmd {
    return nil
}

// Update
func (m DeleteFolderModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch msgType.String() {
        case "esc":
            manager.PopScreen()
        case "enter":
            m.vault.DeleteFolder(m.folder)
            manager.PopScreen()
        }
    }

	return m, nil
}

// View
func (m DeleteFolderModel) View() string {
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

