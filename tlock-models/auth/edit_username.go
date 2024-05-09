package auth

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/context"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/messages"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
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
	usernameError *error
}

// Initialize
func InitializeEditUsernameScreen(user string, context *context.Context) modelmanager.Screen {
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
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, editUsernameKeys.GoBack):
			manager.PopScreen()

		case key.Matches(msgType, editUsernameKeys.Change):
            // New username
			newUsername := screen.newUsername.Value()

            // Rename
            if err := screen.context.Core.RenameUser(screen.user, newUsername); err != nil {
                screen.usernameError = &err
            } else {
                // Pop
                manager.PopScreen()

                // Send updated message
                cmds = append(cmds, func() tea.Msg { return tlockmessages.UserEditedMsg{NewName: newUsername} })
            }
		}
	}

	// Update input box
	screen.newUsername, _ = screen.newUsername.Update(msg)

	// Return
	return screen, tea.Batch(cmds...)
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
