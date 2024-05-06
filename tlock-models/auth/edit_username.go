package auth

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

var editUsernameAscii = `
█▀▀ █▀▄ █ ▀█▀
██▄ █▄▀ █  █`

// Edit username key map
type editUsernameKeyMap struct {
	Change key.Binding
	GoBack key.Binding
}

// ShortHelp()
func (k editUsernameKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Change, k.GoBack}
}

// FullHelp()
func (k editUsernameKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Change},
		{k.GoBack},
	}
}

// Keys
var editUsernameKeys = changePasswordKeyMap{
	Change: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "change username"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

// Edit username screen
type EditUsernameScreen struct {
	// Context
	context *context.Context

	// Username to change
	user string

	// New username
	newUsername textinput.Model

	// Username error
	usernameError *string
}

// Initialize
func InitializeEditUsernameScreen(user string, _ *tlockvault.Vault, context *context.Context) modelmanager.Screen {
	newUsername := components.InitializeInputBox("Your new name goes here...")
	newUsername.Focus()

	return EditUsernameScreen{
		context:     context,
		user:        user,
		newUsername: newUsername,
	}
}

// Init
func (screen EditUsernameScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen EditUsernameScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, editUsernameKeys.GoBack):
			manager.PopScreen()

		case key.Matches(msgType, editUsernameKeys.Change):
			newUsername := screen.newUsername.Value()

			if newUsername == "" {
				screen.usernameError = &USERNAME_EMPTY
			} else {
				if screen.context.Core.Exists(newUsername) {
					screen.usernameError = &USERNAME_EXISTS
				} else {
					screen.context.Core.RenameUser(screen.user, newUsername)

					manager.PopScreen()
				}
			}
		}
	}

	// Update input box
	screen.newUsername, _ = screen.newUsername.Update(msg)

	// Return
	return screen, nil
}

// View
func (screen EditUsernameScreen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(editUsernameAscii), "",
		tlockstyles.Styles.SubText.Render("Change username"), "",
		components.InputGroup("New username", "Choose a new username that you want to set for yourself", screen.usernameError, screen.newUsername),
		tlockstyles.Help.View(editUsernameKeys),
	)
}
