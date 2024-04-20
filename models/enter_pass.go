package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
	tlockcore "github.com/eklairs/tlock/tlock-core"
	tlockvault "github.com/eklairs/tlock/tlock-vault"

	. "github.com/eklairs/tlock/internal/modelmanager"
)

var ENTER_PASS_SIZE = 65

// Enter pass key map
type enterPassKeyMap struct {
    Login key.Binding
    Back key.Binding
}

// ShortHelp()
func (k enterPassKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Login, k.Back}
}

// FullHelp()
func (k enterPassKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
        {k.Login},
        {k.Back},
	}
}

// Enter pass keys
var enterPassKeys = enterPassKeyMap{
	Login: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "login"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

// Enter pass Model
type EnterPassModel struct {
    // Common styles
    styles Styles

    // Password input
    passInput textinput.Model

    // Help
    help help.Model

    // Core
    core tlockcore.TLockCore

    // User spec
    userSpec tlockcore.UserSpec

    // Any error message
    errorMessage bool
}

// Initialize root model
func InitializeEnterPassModel(core tlockcore.TLockCore, userSpec tlockcore.UserSpec) EnterPassModel {
    // Password input
    passwordInput := InitializeInputBox("Your password goes here...")
    passwordInput.Focus()

    return EnterPassModel {
        core: core,
        help: help.New(),
        userSpec: userSpec,
        passInput: passwordInput,
        styles: InitializeStyles(65),
    }
}

// Init
func (m EnterPassModel) Init() tea.Cmd {
    return nil
}

// Update
func (m EnterPassModel) Update(msg tea.Msg, manager *ModelManager) (Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msgType, enterPassKeys.Back):
            manager.PopScreen()
        case key.Matches(msgType, enterPassKeys.Login):
            vault, err := tlockvault.Load(m.userSpec.Vault, m.passInput.Value())

            if err != nil {
                m.errorMessage = true
            } else {
                manager.PushScreen(InitializeDashboardModel(*vault))
            }
        }
    }

    // Update input box
    m.passInput, _ = m.passInput.Update(msg)

	return m, nil
}

// View
func (m EnterPassModel) View() string {
    // List of items
    items := []string {
        m.styles.dimmedCenter.Render(fmt.Sprintf("Login in as %s", m.userSpec.Username)), "",
        m.styles.title.Render("Password"), // Username header
        m.styles.dimmed.Render("Enter the super secret password"), // Username description
        m.styles.input.Render(m.passInput.View()),
    }

    // Add error message if any
    if m.errorMessage {
        items = append(items, "", m.styles.error.Render("Invalid password, please check if it is correct"))
    }

    // Rest
    items = append(items, "", m.styles.center.Render(m.help.View(enterPassKeys)))

    return lipgloss.JoinVertical(
        lipgloss.Left,
        items...
    )
}

