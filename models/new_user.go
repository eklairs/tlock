package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/internal/modelmanager"
)

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
}

// Root Model
type NewUserModel struct {
    styles newUserStyles
    usernameInput textinput.Model
    passwordInput textinput.Model
}

// Initialize root model
func InitializeNewUserModel() NewUserModel {
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
        },
        usernameInput: usernameInput,
        passwordInput: passwordInput,
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
    )
}

