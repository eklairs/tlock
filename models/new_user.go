package models

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/internal/modelmanager"
	tlockcore "github.com/eklairs/tlock/tlock-core"
)

type  newUserKeys struct {
    Tab key.Binding
    Create key.Binding
    GoBack key.Binding
}

func (k  newUserKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Create, k.GoBack}
}

func (k  newUserKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
        {k.Tab},
        {k.Create},
        {k.GoBack},
	}
}

var new_user_keys =  newUserKeys{
	Tab: key.NewBinding(
		key.WithKeys("tab", "shift+tab"),
		key.WithHelp("tab/shift+tab", "switch input"),
	),
	Create: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "create"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

var ascii = `
 _____              _____             
|   | |___ _ _ _   |  |  |___ ___ ___ 
| | | | -_| | | |  |  |  |_ -| -_|  _|
|_|___|___|_____|  |_____|___|___|_|  
`

type newUserStyles struct {
    title lipgloss.Style
    titleCenter lipgloss.Style
    dimmed lipgloss.Style
    input lipgloss.Style
    center lipgloss.Style
}

// Root Model
type NewUserModel struct {
    styles newUserStyles
    help help.Model
    key newUserKeys
    core tlockcore.TLockCore
    usernameInput textinput.Model
    passwordInput textinput.Model
}

// Initialize root model
func InitializeNewUserModel(core tlockcore.TLockCore) NewUserModel {
    dimmed := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

    usernameInput := textinput.New();
    usernameInput.Prompt = ""
    usernameInput.Width = 58
    usernameInput.Placeholder = "Your username goes here..."
    usernameInput.PlaceholderStyle = dimmed.Copy().Background(lipgloss.Color("#1e1e2e"))
    usernameInput.Focus() // Let us have the username input focused at first

    passwordInput := textinput.New();
    passwordInput.Prompt = ""
    passwordInput.Width = 58
    passwordInput.Placeholder = "Your password goes here..."
    passwordInput.PlaceholderStyle = dimmed.Copy().Background(lipgloss.Color("#1e1e2e"))

    return NewUserModel {
        styles: newUserStyles{
            title: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")),
            titleCenter: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")).Width(65).Align(lipgloss.Center),
            input: lipgloss.NewStyle().Padding(1, 3).Width(65).Background(lipgloss.Color("#1e1e2e")),
            dimmed: dimmed,
            center: lipgloss.NewStyle().Align(lipgloss.Center).Width(65),
        },
        usernameInput: usernameInput,
        passwordInput: passwordInput,
        core: core,
        help: help.New(),
        key: new_user_keys,
    }
}

// Init
func (m NewUserModel) Init() tea.Cmd {
    return nil
}

// Update
func (m NewUserModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch msgType.String() {
        case "tab":
            if m.usernameInput.Focused() {
                m.usernameInput.Blur()
                m.passwordInput.Focus()
            } else {
                m.usernameInput.Focus()
                m.passwordInput.Blur()
            }
        case "esc":
            manager.PopScreen()
        case "enter":
            vault := m.core.Users.AddNewUser(m.usernameInput.Value(), m.passwordInput.Value())

            manager.PushScreen(InitializeDashboardModel(vault))
        }
    }

    if m.usernameInput.Focused() {
        m.usernameInput, _ = m.usernameInput.Update(msg)
    }

    if m.passwordInput.Focused() {
        m.passwordInput, _ = m.passwordInput.Update(msg)
    }

	return m, nil
}

// View
func (m NewUserModel) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Left,
        m.styles.titleCenter.Render(ascii), // Title
        m.styles.title.Render("Username"), // Username header
        m.styles.dimmed.Render("Choose an awesome name, how about Komaru?"), // Username description 
        m.styles.input.Render(m.usernameInput.View()), "",
        m.styles.title.Render("Password"), // Username header
        m.styles.dimmed.Render("Shush! This is a super secret password"), // Username description
        m.styles.input.Render(m.passwordInput.View()), "",
        m.styles.center.Render(m.help.View(new_user_keys)),
    )
}
