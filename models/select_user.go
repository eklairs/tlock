package models

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/internal/modelmanager"
	tlockcore "github.com/eklairs/tlock/tlock-core"
)

var _ascii = `
 _____     _         _      _____             
|   __|___| |___ ___| |_   |  |  |___ ___ ___ 
|__   | -_| | -_|  _|  _|  |  |  |_ -| -_|  _|
|_____|___|_|___|___|_|    |_____|___|___|_|  
`

type selectUserStyles struct {
    title lipgloss.Style
    titleCenter lipgloss.Style
    dimmed lipgloss.Style
    dimmedCenter lipgloss.Style
    input lipgloss.Style
    userItem lipgloss.Style
    userItemFocused lipgloss.Style
}

// Root Model
type SelectUserModel struct {
    styles selectUserStyles
    focused_index int
    core tlockcore.TLockCore
}

// Initialize root model
func InitializeSelectUserModel(core tlockcore.TLockCore) SelectUserModel {
    dimmed := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

    return SelectUserModel {
        styles: selectUserStyles{
            title: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")),
            titleCenter: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")).Width(65).Align(lipgloss.Center),
            input: lipgloss.NewStyle().Padding(1, 3).Width(65).Background(lipgloss.Color("#1e1e2e")),
            dimmed: dimmed,
            dimmedCenter: dimmed.Width(65).Copy().Align(lipgloss.Center),
            userItem: lipgloss.NewStyle().Padding(1, 3).Width(65),
            userItemFocused: lipgloss.NewStyle().Padding(1, 3).Width(65).Background(lipgloss.Color("#1E1E2E")).Foreground(lipgloss.Color("12")).Bold(true),
        },
        core: core,
    }
}

// Init
func (m SelectUserModel) Init() tea.Cmd {
    return nil
}

// Update
func (m SelectUserModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch msgType.String() {
        case tea.KeyDown.String(), "j":
            total_users := len(m.core.Users.Users)

            m.focused_index = (m.focused_index + 1) % total_users
        
        case tea.KeyUp.String(), "k":
            total_users := len(m.core.Users.Users)

            if m.focused_index == 0 {
                m.focused_index = total_users - 1
            } else {
                m.focused_index -= 1
            }
        }
    }

	return m, nil
}

// View
func (m SelectUserModel) View() string {
    user_items := []string {
        m.styles.titleCenter.Render(_ascii), // Title
        m.styles.dimmedCenter.Render("Select a user to continue as"), "",
    }

    for index, user := range m.core.Users.Users {
        render_fn := m.styles.userItem.Render

        if index == m.focused_index {
            render_fn = m.styles.userItemFocused.Render
        }

        renderable := render_fn(
            lipgloss.JoinHorizontal(
                lipgloss.Center,
                user.Username,
                strings.Repeat(" ", 65 - len(user.Username) - 1 - 6),
                "›",
            ),
        )

        user_items = append(user_items, renderable)
    }

    return lipgloss.JoinVertical(
        lipgloss.Left,
        user_items...,
    )
}

