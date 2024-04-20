package models

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
	tlockcore "github.com/eklairs/tlock/tlock-core"

	. "github.com/eklairs/tlock/internal/modelmanager"
)

var CREATE_USER_SIZE = 65

// Create user model key bindings
type createUserKeyMap struct {
    Tab key.Binding
    Create key.Binding
    GoBack key.Binding
}

// ShortHelp()
func (k createUserKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Create, k.GoBack}
}

// FullHelp()
func (k createUserKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
        {k.Tab},
        {k.Create},
        {k.GoBack},
	}
}

// Keys
var createUserKeys =  createUserKeyMap{
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

// Create user model
type NewUserModel struct {
    // Help
    help help.Model

    // Common styles
    styles Styles

    // Core
    core tlockcore.TLockCore

    // Username input
    usernameInput textinput.Model

    // Password input
    passwordInput textinput.Model
}

// Initialize new user model
func InitializeNewUserModel(core tlockcore.TLockCore) NewUserModel {
    // Username input box
    usernameInput := InitializeInputBox("Your username goes here...")
    usernameInput.Focus() // Let us have the username input focused at first

    // Password input
    passwordInput := InitializeInputBox("Your password goes here...");

    // Return
    return NewUserModel {
        core: core,
        help: help.New(),
        usernameInput: usernameInput,
        passwordInput: passwordInput,
        styles: InitializeStyles(CREATE_USER_SIZE),
    }
}

// Init
func (m NewUserModel) Init() tea.Cmd {
    return nil
}

// Update
func (m NewUserModel) Update(msg tea.Msg, manager *ModelManager) (Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msgType, createUserKeys.Tab):
            if m.usernameInput.Focused() {
                m.usernameInput.Blur()
                m.passwordInput.Focus()
            } else {
                m.usernameInput.Focus()
                m.passwordInput.Blur()
            }

        case key.Matches(msgType, createUserKeys.GoBack):
            manager.PopScreen()

        case key.Matches(msgType, createUserKeys.Create):
            m.core.Users.AddNewUser(m.usernameInput.Value(), m.passwordInput.Value())
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
        m.styles.title.Render("Username"), // Username header
        m.styles.dimmed.Render("Choose an awesome name, how about Komaru?"), // Username description 
        m.styles.input.Render(m.usernameInput.View()), "",
        m.styles.title.Render("Password"), // Username header
        m.styles.dimmed.Render("Shush! This is a super secret password"), // Username description
        m.styles.input.Render(m.passwordInput.View()), "",
        m.styles.center.Render(m.help.View(createUserKeys)),
    )
}

