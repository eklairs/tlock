package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/internal/boundedinteger"
	"github.com/eklairs/tlock/internal/modelmanager"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

var ______ascii = `
 _____   _ _ _   
|   __|_| |_| |_ 
|   __| . | |  _|
|_____|___|_|_|  
`

type moveTokenStyles struct {
    title lipgloss.Style
    titleCenter lipgloss.Style
    dimmed lipgloss.Style
    dimmedCenter lipgloss.Style
    input lipgloss.Style
    userItem lipgloss.Style
    userItemFocused lipgloss.Style
}

// Root Model
type MoveTokenModel struct {
    styles moveTokenStyles
    vault *tlockvault.TLockVault
    folder int
    token int
    focused_index boundedinteger.BoundedInteger
}

// Initialize root model
func InitializeMoveTokenModel(vault *tlockvault.TLockVault, folder, token int) MoveTokenModel {
    dimmed := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
    styles := moveTokenStyles{
        title: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")),
        titleCenter: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")).Width(55).Align(lipgloss.Center),
        input: lipgloss.NewStyle().Padding(1, 3).Width(55).Background(lipgloss.Color("#1e1e2e")),
        dimmed: dimmed,
        dimmedCenter: dimmed.Width(55).Copy().Align(lipgloss.Center),
        userItem: lipgloss.NewStyle().Padding(1, 3).Width(65).Foreground(lipgloss.Color("8")),
        userItemFocused: lipgloss.NewStyle().Padding(1, 3).Width(65).Background(lipgloss.Color("#1E1E2E")).Foreground(lipgloss.Color("12")),
    }

    return MoveTokenModel {
        styles: styles,
        vault: vault,
        folder: folder,
        token: token,
        focused_index: boundedinteger.New(0, len(vault.Data.Folders)),
    }
}

// Init
func (m MoveTokenModel) Init() tea.Cmd {
    return nil
}

// Update
func (m MoveTokenModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch msgType.String() {
        case tea.KeyDown.String(), "j":
            m.focused_index.Increase()   
        case tea.KeyUp.String(), "k":
            m.focused_index.Decrease()
        case "esc":
            manager.PopScreen()
        case "enter":
            m.vault.MoveURI(m.folder, m.token, m.focused_index.Value)
            manager.PopScreen()
        }
    }

	return m, nil
}

// View
func (m MoveTokenModel) View() string {
    user_items := []string {
        m.styles.titleCenter.Render(______ascii), // Title
        m.styles.dimmedCenter.Render("Select the folder to move the token to"), "",
    }

    for index, folder := range m.vault.Data.Folders {
        render_fn := m.styles.userItem.Render

        if index == m.focused_index.Value {
            render_fn = m.styles.userItemFocused.Render
        }

        renderable := render_fn(
            lipgloss.JoinHorizontal(
                lipgloss.Center,
                folder.Name,
                strings.Repeat(" ", 65 - len(folder.Name) - 8 - 6),
                fmt.Sprintf("%d tokens", len(folder.Uris)),
            ),
        )

        user_items = append(user_items, renderable)
    }

    return lipgloss.JoinVertical(
        lipgloss.Left,
        user_items...
    )
}

