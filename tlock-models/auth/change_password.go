package auth

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tlockcore "github.com/eklairs/tlock/tlock-core"
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
	Tab    key.Binding
	Change key.Binding
	GoBack key.Binding
}

// ShortHelp()
func (k changePasswordKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Change, k.GoBack}
}

// FullHelp()
func (k changePasswordKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab},
		{k.Change},
		{k.GoBack},
	}
}

// Keys
var changePasswordKeys = changePasswordKeyMap{
	Tab: key.NewBinding(
		key.WithKeys("tab", "shift+tab"),
		key.WithHelp("tab/shift+tab", "switch input"),
	),
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

	// Original password input
	originalPassword textinput.Model

	// New password input
	newPassword textinput.Model

	// Original password error
	passwordError *string

	// Vault instance
	// If it is nil, this means that its locked with password
	vault *tlockvault.Vault

	// User
	user tlockcore.User
}

// Initializes a new instance of the create user screen
func InitializeChangePasswordScreen(context *context.Context, user tlockcore.User) ChangePasswordScreen {
	// Try to load the vault
	vault, _ := tlockvault.Load(user.Vault, "")

	// Input box for original password
	originalPassword := components.InitializeInputBox("Your original password goes here...")
	originalPassword.EchoMode = textinput.EchoPassword
	originalPassword.EchoCharacter = constants.CHAR_ECHO
	originalPassword.Focus()

	// Input box for password
	newPassword := components.InitializeInputBox("Your new password goes here...")
	newPassword.EchoMode = textinput.EchoPassword
	newPassword.EchoCharacter = constants.CHAR_ECHO

	// If the vault is unlocked, focus the new password input box
	if vault != nil {
		originalPassword.Blur()
		newPassword.Focus()
	}

	return ChangePasswordScreen{
		context:          context,
		originalPassword: originalPassword,
		newPassword:      newPassword,
		user:             user,
		passwordError:    nil,
		vault:            vault,
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
		screen.passwordError = nil

		switch {
		case key.Matches(msgType, changePasswordKeys.GoBack):
			manager.PopScreen()

		case key.Matches(msgType, changePasswordKeys.Tab):
			// If the vault is loaded, there is only one input
			// Then skip
			if screen.vault != nil {
				break
			}

			// Switch focus
			if screen.originalPassword.Focused() {
				screen.originalPassword.Blur()
				screen.newPassword.Focus()
			} else {
				screen.originalPassword.Focus()
				screen.newPassword.Blur()
			}

		case key.Matches(msgType, changePasswordKeys.Change):
			newPassword := screen.newPassword.Value()

			// If the vault is already open, lets just change the password and get away
			if screen.vault != nil {
				screen.vault.ChangePassword(newPassword)
			} else {
				// Load vault with the given password
				vault, err := tlockvault.Load(screen.user.Vault, screen.originalPassword.Value())

				// Check if there is any error
				if err != nil {
					screen.passwordError = &ERROR_PASSWORD_WRONG
				} else {
					// Update password
					vault.ChangePassword(newPassword)

					// Pop screen
					manager.PopScreen()
				}
			}

		default:
			// Update input boxes
			if screen.originalPassword.Focused() {
				screen.originalPassword, _ = screen.originalPassword.Update(msg)
			}

			if screen.newPassword.Focused() {
				screen.newPassword, _ = screen.newPassword.Update(msg)
			}
		}
	}

	return screen, cmd
}

// View
func (screen ChangePasswordScreen) View() string {
	items := []string{
		tlockstyles.Styles.Title.Render(changePasswordAsciiArt), "",
		tlockstyles.Styles.SubText.Render("Change your password"), "",
	}

	// Add original password requirement if the vault is locked
	if screen.vault == nil {
		items = append(items, components.InputGroup("Original password", "Enter your original password, the one that you used to use to login", screen.passwordError, screen.originalPassword))
	}

	// Add Rest
	items = append(
		items,
		components.InputGroup("New password", "Enter the new password that you want to use to login from next time", nil, screen.newPassword),
		tlockstyles.Help.View(changePasswordKeys),
	)

	// Return
	return lipgloss.JoinVertical(lipgloss.Center, items...)
}
