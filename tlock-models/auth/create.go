package auth

import (
	"os/user"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/constants"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-models/dashboard"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

// Errors
var USERNAME_EXISTS = "User with that name already exists"
var USERNAME_EMPTY = "Please enter a username"

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
	context *context.Context

	// Username input
	usernameInput textinput.Model

	// Password input
	passwordInput textinput.Model

	// Username error message
	usernameError *string

	// System user if found
	systemUser *user.User
}

// Initializes a new instance of the create user screen
func InitializeCreateUserScreen(context *context.Context) CreateUserScreen {
	// Placeholder
	placeholder := "Your username goes here..."

	// Get current user
	user, err := user.Current()

	if err == nil {
		placeholder = user.Username
	} else {
		placeholder = "Cannot find your username, please enter one..."
	}

	// Input box for username
	usernameInput := components.InitializeInputBox(placeholder)
	usernameInput.Focus()

	// Input box for password
	passwordInput := components.InitializeInputBox("Your password goes here...")
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = constants.CHAR_ECHO

	return CreateUserScreen{
		context:       context,
		usernameInput: usernameInput,
		passwordInput: passwordInput,
		systemUser:    user,
	}
}

// Init
func (screen CreateUserScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen CreateUserScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	var cmd tea.Cmd

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
			username := screen.usernameInput.Value()

			if username == "" {
				if screen.systemUser != nil {
					username = screen.systemUser.Username
				} else {
					// Set error
					screen.usernameError = &USERNAME_EMPTY

					// Break
					break
				}
			}

			// Add new user
			vault, err := screen.context.Core.AddNewUser(username, screen.passwordInput.Value())

			if err != nil {
				screen.usernameError = &USERNAME_EXISTS
			} else {
				// Push dashboard screen
				cmd = manager.PushScreen(dashboard.InitializeDashboardScreen(screen.usernameInput.Value(), vault, screen.context))
			}

		default:
			// Update input boxes
			if screen.usernameInput.Focused() {
				screen.usernameInput, _ = screen.usernameInput.Update(msg)
			}

			if screen.passwordInput.Focused() {
				screen.passwordInput, _ = screen.passwordInput.Update(msg)
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
		components.InputGroup("Username", "Choose an awesome username, or keep it empty to use the current system name", screen.usernameError, screen.usernameInput),
		components.InputGroup("Password", "Choose a super strong password, or keep it empty if you don't want any password", nil, screen.passwordInput),
		tlockstyles.Help.View(createUserKeys),
	)
}
