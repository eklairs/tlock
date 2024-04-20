package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/internal/modelmanager"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

var ___ascii = `
 _____              _____     _   _         
|   | |___ _ _ _   |   __|___| |_| |___ ___ 
| | | | -_| | | |  |   __| . | | . | -_|  _|
|_|___|___|_____|  |__|  |___|_|___|___|_|  
`

type newFolderStyles struct {
    title lipgloss.Style
    titleCenter lipgloss.Style
    dimmed lipgloss.Style
    dimmedCenter lipgloss.Style
    input lipgloss.Style
}

// Root Model
type NewFolderModel struct {
    styles newFolderStyles
    name textinput.Model
    vault *tlockvault.TLockVault
}

// Initialize root model
func InitializeNewFolderModel(vault *tlockvault.TLockVault) NewFolderModel {
    dimmed := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

    name := textinput.New();
    name.Prompt = ""
    name.Width = 58
    name.Placeholder = "Your password goes here..."
    name.PlaceholderStyle = dimmed.Copy().Background(lipgloss.Color("#1e1e2e"))
    name.Focus()

    return NewFolderModel {
        styles: newFolderStyles{
            title: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")),
            titleCenter: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")).Width(65).Align(lipgloss.Center),
            input: lipgloss.NewStyle().Padding(1, 3).Width(65).Background(lipgloss.Color("#1e1e2e")),
            dimmed: dimmed,
            dimmedCenter: dimmed.Width(65).Copy().Align(lipgloss.Center),
        },
        vault: vault,
        name: name,
    }
}

// Init
func (m NewFolderModel) Init() tea.Cmd {
    return nil
}

// Update
func (m NewFolderModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch msgType.String() {
        case "enter":
            m.vault.AddFolder(m.name.Value())
            manager.PopScreen()
        }
    }

    m.name, _ = m.name.Update(msg)

	return m, nil
}

// View
func (m NewFolderModel) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Left,
        m.styles.titleCenter.Render(___ascii), "", // Title
        m.styles.title.Render("Name"), // Username header
        m.styles.dimmed.Render("Choose a name for your folder, like Socials!"), // Username description
        m.styles.input.Render(m.name.View()), "",
    )
}

