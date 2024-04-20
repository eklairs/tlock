package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/internal/modelmanager"
	tlockcore "github.com/eklairs/tlock/tlock-core"
)

var __ascii = `
 _____     _              _____             
|   __|___| |_ ___ ___   |  _  |___ ___ ___ 
|   __|   |  _| -_|  _|  |   __| .'|_ -|_ -|
|_____|_|_|_| |___|_|    |__|  |__,|___|___|
`

type enterPassStyles struct {
    title lipgloss.Style
    titleCenter lipgloss.Style
    dimmed lipgloss.Style
    dimmedCenter lipgloss.Style
    input lipgloss.Style
}

// Root Model
type EnterPassModel struct {
    styles enterPassStyles
    passInput textinput.Model
    core tlockcore.TLockCore
    userSpec tlockcore.UserSpec
}

// Initialize root model
func InitializeEnterPassModel(core tlockcore.TLockCore, userIndex int) EnterPassModel {
    dimmed := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

    passwordInput := textinput.New();
    passwordInput.Prompt = ""
    passwordInput.Width = 58
    passwordInput.Placeholder = "Your password goes here..."
    passwordInput.PlaceholderStyle = dimmed.Copy().Background(lipgloss.Color("#1e1e2e"))
    passwordInput.Focus()

    return EnterPassModel {
        styles: enterPassStyles{
            title: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")),
            titleCenter: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")).Width(65).Align(lipgloss.Center),
            input: lipgloss.NewStyle().Padding(1, 3).Width(65).Background(lipgloss.Color("#1e1e2e")),
            dimmed: dimmed,
            dimmedCenter: dimmed.Width(65).Copy().Align(lipgloss.Center),
        },
        core: core,
        passInput: passwordInput,
        userSpec: core.Users.Users[userIndex],
    }
}

// Init
func (m EnterPassModel) Init() tea.Cmd {
    return nil
}

// Update
func (m EnterPassModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch msgType.String() {

        }
    }

    m.passInput, _ = m.passInput.Update(msg)

	return m, nil
}

// View
func (m EnterPassModel) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Left,
        m.styles.titleCenter.Render(__ascii), // Title
        m.styles.dimmedCenter.Render(fmt.Sprintf("Login in as %s", m.userSpec.Username)), "",
        m.styles.title.Render("Password"), // Username header
        m.styles.dimmed.Render("Enter the super secret password"), // Username description
        m.styles.input.Render(m.passInput.View()), "",
    )
}
