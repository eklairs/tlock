package auth

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/constants"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"

	tlockcore "github.com/eklairs/tlock/tlock-core"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	tlockstyles "github.com/eklairs/tlock/tlock/styles"
)

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

// Next function
type NextFunc = func(string, *tlockvault.Vault, *context.Context) modelmanager.Screen

// Enter pass screen
type EnterPassScreen struct {
	// Context
	context *context.Context

	// Password input
	passInput textinput.Model

	// User spec
	user tlockcore.User

	// Any error message
	errorMessage *error

	// Next
	next NextFunc

	// Ascii
	ascii string

	// Description
	description string
}

// Initialize root model
func InitializeEnterPassScreen(context *context.Context, user tlockcore.User, next NextFunc) EnterPassScreen {
	return InitializeEnterPassScreenCustomOpts(context, user, next, enterPassAsciiArt, "Login in as %s")
}

// Initializes enter pass screen with custom title and desc
func InitializeEnterPassScreenCustomOpts(context *context.Context, user tlockcore.User, next NextFunc, ascii, desc string) EnterPassScreen {
	// Password input
	passwordInput := components.InitializeInputBox("Your password goes here...")
	passwordInput.EchoCharacter = constants.CHAR_ECHO
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.Focus()

	return EnterPassScreen{
		context:     context,
		user:        user,
		passInput:   passwordInput,
		next:        next,
		ascii:       ascii,
		description: desc,
	}
}

// Init
func (screen EnterPassScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen EnterPassScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	var cmd tea.Cmd

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		if screen.passInput.Value() != "" {
			screen.errorMessage = nil
		}

		switch {
		case strings.Contains(msgType.String(), "tab"):
			// We dont want to allow tabs!
		case key.Matches(msgType, enterPassKeys.Back):
			manager.PopScreen()
		case key.Matches(msgType, enterPassKeys.Login):
			vault, err := tlockvault.Load(screen.user.Vault(), screen.passInput.Value())

			// Show error message if vault was failed to be unlocked
			if err != nil {
				screen.errorMessage = &err
			} else {
				cmd = manager.ReplaceScreen(screen.next(screen.user.S(), vault, screen.context))
			}
		default:
			// Update input box
			screen.passInput, _ = screen.passInput.Update(msg)
		}
	}

	return screen, cmd
}

// View
func (screen EnterPassScreen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		tlockstyles.Title(screen.ascii), "",
		tlockstyles.Dimmed(fmt.Sprintf(screen.description, screen.user.S())), "",
		components.InputGroup("Password", "Enter the super secret password", screen.errorMessage, screen.passInput),
		tlockstyles.HelpView(enterPassKeys),
	)
}
