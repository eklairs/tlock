package auth

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tlockcore "github.com/eklairs/tlock/tlock-core"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/constants"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-models/dashboard"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

var ERROR_PASSWORD_WRONG = "Wrong password, please check if it is correct"

var enterPassAsciiArt = `
█▀█ ▄▀█ █▀ █▀ █ █ █ █▀█ █▀█ █▀▄
█▀▀ █▀█ ▄█ ▄█ ▀▄▀▄▀ █▄█ █▀▄ █▄▀`

// Enter pass key map
type enterPassKeyMap struct {
	Login key.Binding
	Back  key.Binding
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

// Enter pass screen
type EnterPassScreen struct {
	// Context
	context context.Context

	// Password input
	passInput textinput.Model

	// User spec
	user tlockcore.User

	// Any error message
	errorMessage *string
}

// Initialize root model
func InitializeEnterPassScreen(context context.Context, user tlockcore.User) EnterPassScreen {
	// Password input
	passwordInput := components.InitializeInputBox("Your password goes here...")
	passwordInput.EchoCharacter = constants.CHAR_ECHO
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.Focus()

	return EnterPassScreen{
		context:   context,
		user:      user,
		passInput: passwordInput,
	}
}

// Init
func (screen EnterPassScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen EnterPassScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	var cmd tea.Cmd

	// Update input box
	screen.passInput, _ = screen.passInput.Update(msg)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		if screen.passInput.Value() != "" {
			screen.errorMessage = nil
		}

		switch {
		case key.Matches(msgType, enterPassKeys.Back):
			manager.PopScreen()
		case key.Matches(msgType, enterPassKeys.Login):
			vault, err := tlockvault.Load(screen.user.Vault, screen.passInput.Value())

			// Show error message if vault was failed to be unlocked
			if err != nil {
				screen.errorMessage = &ERROR_PASSWORD_WRONG
			} else {
				cmd = manager.PushScreen(dashboard.InitializeDashboardScreen(*vault, screen.context))
			}
		}
	}

	return screen, cmd
}

// View
func (screen EnterPassScreen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(enterPassAsciiArt), "",
		tlockstyles.Styles.SubText.Render(fmt.Sprintf("Login in as %s", screen.user.Username)), "",
		components.InputGroup("Password", "Enter the super secret password", screen.errorMessage, screen.passInput),
		tlockstyles.Help.View(enterPassKeys),
	)
}
