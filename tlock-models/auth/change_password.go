package auth

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/constants"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

// Change password ascii art
var changePasswordAsciiArt = `
█▀▀ █▀█ █▀▀ ▄▀█ ▀█▀ █▀▀   █ █ █▀ █▀▀ █▀█
█▄▄ █▀▄ ██▄ █▀█  █  ██▄   █▄█ ▄█ ██▄ █▀▄`

// Change password key map
type changePasswordKeyMap struct {
	Change key.Binding
	GoBack key.Binding
}

// ShortHelp()
func (k changePasswordKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Change, k.GoBack}
}

// FullHelp()
func (k changePasswordKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Change},
		{k.GoBack},
	}
}

// Keys
var changePasswordKeys = changePasswordKeyMap{
	Change: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "change password"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

// Change password user screen
type ChangePasswordScreen struct {
	// Context
	context *context.Context

	// New password input
	newPassword textinput.Model

	// Vault instance
	vault *tlockvault.Vault

	// User
	user string
}

// Initializes a new instance of the create user screen
func InitializeChangePasswordScreen(context *context.Context, vault *tlockvault.Vault, user string) ChangePasswordScreen {
	// Input box for password
	newPassword := components.InitializeInputBox("Your new password goes here...")
	newPassword.EchoMode = textinput.EchoPassword
	newPassword.EchoCharacter = constants.CHAR_ECHO

	return ChangePasswordScreen{
		context:     context,
		newPassword: newPassword,
		user:        user,
		vault:       vault,
	}
}

// Init
func (screen ChangePasswordScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen ChangePasswordScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	var cmd tea.Cmd

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, changePasswordKeys.GoBack):
			manager.PopScreen()

		case key.Matches(msgType, changePasswordKeys.Change):
			// Change password
			screen.vault.ChangePassword(screen.newPassword.Value())

		default:
			screen.newPassword, _ = screen.newPassword.Update(msg)
		}
	}

	return screen, cmd
}

// View
func (screen ChangePasswordScreen) View() string {
	items := []string{
		tlockstyles.Styles.Title.Render(changePasswordAsciiArt), "",
		tlockstyles.Styles.SubText.Render("Change your password"), "",
		components.InputGroup("New password", "Enter the new password that you want to use to login from next time", nil, screen.newPassword),
		tlockstyles.Help.View(changePasswordKeys),
	}

	// Return
	return lipgloss.JoinVertical(lipgloss.Center, items...)
}
