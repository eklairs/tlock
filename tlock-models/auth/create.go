package auth

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-models/dashboard"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

// Errors
var USERNAME_EMPTY = "Please choose a username"

// Create user ascii art
var createUserAsciiArt = `
█▀▀ █▀█ █▀▀ ▄▀█ ▀█▀ █▀▀   █ █ █▀ █▀▀ █▀█
█▄▄ █▀▄ ██▄ █▀█  █  ██▄   █▄█ ▄█ ██▄ █▀▄`

// Create user key map
type createUserKeyMap struct {
	Tab    key.Binding
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
var createUserKeys = createUserKeyMap{
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

// Create user screen
type CreateUserScreen struct {
	// Context
	context context.Context

	// Username input
	usernameInput textinput.Model

	// Password input
	passwordInput textinput.Model

	// Username error message
	usernameError *string
}

// Initializes a new instance of the create user screen
func InitializeCreateUserScreen(context context.Context) CreateUserScreen {
	// Input box for username
	usernameInput := components.InitializeInputBox("Your username goes here...")
	usernameInput.Focus()

	return CreateUserScreen{
		context:       context,
		usernameInput: usernameInput,
		passwordInput: components.InitializeInputBox("Your password goes here..."),
	}
}

// Init
func (screen CreateUserScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen CreateUserScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	var cmd tea.Cmd

	// Update input boxes
	if screen.usernameInput.Focused() {
		screen.usernameInput, _ = screen.usernameInput.Update(msg)
	}

	if screen.passwordInput.Focused() {
		screen.passwordInput, _ = screen.passwordInput.Update(msg)
	}

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		// Remove error (if any) if the user input box is not empty
		if screen.usernameInput.Value() != "" {
			screen.usernameError = nil
		}

		switch {
		case key.Matches(msgType, createUserKeys.GoBack):
			manager.PopScreen()

		case key.Matches(msgType, createUserKeys.Tab):
			if screen.usernameInput.Focused() {
				screen.usernameInput.Blur()
				screen.passwordInput.Focus()
			} else {
				screen.usernameInput.Focus()
				screen.passwordInput.Blur()
			}

		case key.Matches(msgType, createUserKeys.Create):
			if screen.usernameInput.Value() == "" {
				screen.usernameError = &USERNAME_EMPTY
			} else {
				// Add new user
				vault := screen.context.Core.AddNewUser(screen.usernameInput.Value(), screen.passwordInput.Value())

				// Push dashboard screen
				cmd = manager.PushScreen(dashboard.InitializeDashboardScreen(vault, screen.context))
			}
		}
	}

	return screen, cmd
}

// View
func (screen CreateUserScreen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(createUserAsciiArt), "",
		tlockstyles.Styles.SubText.Render("Create a new user"), "",
		components.InputGroup("Username", "Choose an awesome username, like Komaru!", screen.usernameError, screen.usernameInput),
		components.InputGroup("Password", "Choose a super strong password, or keep it empty if you don't want any password", nil, screen.passwordInput),
		tlockstyles.Help.View(createUserKeys),
	)
}
