package auth

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/context"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/messages"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock/styles"
)

var deleteUserAsciiArt = `
█▀▄ █▀▀ █   █▀▀ ▀█▀ █▀▀
█▄▀ ██▄ █▄▄ ██▄  █  ██▄`

// Delete user model key bindings
type deleteUserKeyMap struct {
	Delete key.Binding
	GoBack key.Binding
}

// ShortHelp()
func (k deleteUserKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.GoBack, k.Delete}
}

// FullHelp()
func (k deleteUserKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.GoBack},
		{k.Delete},
	}
}

// Keys
var deleteUserKeys = deleteUserKeyMap{
	Delete: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "delete user"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

// Delete user screen
type DeleteUserScreen struct {
	// User to delete
	User string

	// Context
	Context *context.Context
}

func InitializeDeleteUserScreen(user string, context *context.Context) modelmanager.Screen {
	return DeleteUserScreen{
		User:    user,
		Context: context,
	}
}

// Init
func (screen DeleteUserScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen DeleteUserScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, deleteUserKeys.Delete):
			// Delete user
			screen.Context.Core.DeleteUser(screen.User)

			// Pop
			manager.PopScreen()

			// Send user deleted message
			cmds = append(cmds, func() tea.Msg { return tlockmessages.UserDeletedMsg{} })

		case key.Matches(msgType, deleteUserKeys.GoBack):
			manager.PopScreen()
		}
	}

	return screen, tea.Batch(cmds...)
}

// View
func (screen DeleteUserScreen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(deleteUserAsciiArt), "",
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			tlockstyles.Styles.Base.Render(fmt.Sprintf("Are you sure you want to ")),
			tlockstyles.Styles.Error.Copy().Bold(true).Render(fmt.Sprintf("× DELETE ")),
			tlockstyles.Styles.Base.Render(fmt.Sprintf("user %s forever?", screen.User)),
		), "",
		tlockstyles.Help.View(deleteUserKeys),
	)
}
