package auth

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/context"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/messages"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	tlockstyles "github.com/eklairs/tlock/tlock/styles"
)

const userOptionsAscii = `
█▀█ █▀█ ▀█▀ █ █▀█ █▄ █ █▀
█▄█ █▀▀  █  █ █▄█ █ ▀█ ▄█`

// User options key map
type userOptionsKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Esc   key.Binding
}

// ShortHelp()
func (k userOptionsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Esc, k.Enter}
}

// LongHelp()
func (k userOptionsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up},
		{k.Down},
		{k.Esc},
		{k.Enter},
	}
}

// Keys
var userOptionsKeys = userOptionsKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "choose option"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

// User options screen
type UserOptionsScreen struct {
	// Context
	context *context.Context

	// User
	user string

	// Focused
	focused int

	// Vault
	vault *tlockvault.Vault
}

// Initializes user options screen
func InitializeUserOptionsScreen(user string, vault *tlockvault.Vault, context *context.Context) modelmanager.Screen {
	return UserOptionsScreen{
		context: context,
		user:    user,
		focused: 0,
		vault:   vault,
	}
}

// Init
func (screen UserOptionsScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen UserOptionsScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	cmd := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, userOptionsKeys.Down):
			if screen.focused != 2 {
				screen.focused += 1
			}

		case key.Matches(msgType, userOptionsKeys.Up):
			if screen.focused != 0 {
				screen.focused -= 1
			}

		case key.Matches(msgType, userOptionsKeys.Esc):
			manager.PopScreen()

		case key.Matches(msgType, userOptionsKeys.Enter):
			switch screen.focused {
			case 0:
				cmd = append(cmd, manager.PushScreen(InitializeEditUsernameScreen(screen.user, screen.context)))

			case 1:
				cmd = append(cmd, manager.PushScreen(InitializeChangePasswordScreen(screen.context, screen.vault, screen.user)))

			case 2:
				cmd = append(cmd, manager.PushScreen(InitializeDeleteUserScreen(screen.user, screen.context)))
			}
		}

	case tlockmessages.UserDeletedMsg:
		// If we recieve delete user message, then its time we pop ourself (it is no longer needed / workable)
		manager.PopScreen()

		// Now we are at the select user screen
		// If there are no users, we will replace the screen with create user screen
		if len(screen.context.Core.Users) == 0 {
			manager.ReplaceScreen(InitializeCreateUserScreen(screen.context))
		}

	case tlockmessages.UserEditedMsg:
		// Lets update the name with the new one
		// To prevent any crashes
		screen.user = msgType.NewName
	}

	return screen, tea.Batch(cmd...)
}

// View
func (screen UserOptionsScreen) View() string {
	// Elements
	elements := []string{
		tlockstyles.Styles.Title.Render(userOptionsAscii), "",
		tlockstyles.Styles.SubText.Render(fmt.Sprintf("Select an option for %s", screen.user)), "",
	}

	// Options
	options := []string{"Edit username", "Change password", "Delete"}

	// Render!
	for index, option := range options {
		// Decide the renderer based on focused index
		renderer := components.ListItemInactive

		if index == screen.focused {
			renderer = components.ListItemActive
		}

		// Render
		elements = append(elements, renderer(65, option, "›"))
	}

	// Add help
	elements = append(elements, "", tlockstyles.Help.View(userOptionsKeys))

	return lipgloss.JoinVertical(lipgloss.Center, elements...)
}
