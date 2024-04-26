package auth

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

var selectUserAscii = `
█   █▀█ █▀▀ █ █▄ █
█▄▄ █▄█ █▄█ █ █ ▀█`

// Select user key map
type selectUserKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	New   key.Binding
}

// ShortHelp()
func (k selectUserKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.New, k.Enter}
}

// LongHelp()
func (k selectUserKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up},
		{k.Down},
		{k.New},
		{k.Enter},
	}
}

// Keys
var selectUserKeys = selectUserKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	New: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "new user"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "login as"),
	),
}

// Select user
type SelectUserScreen struct {
	// Context
	context context.Context

	// List view
	listview list.Model
}

// New instance of select user
func InitializeSelectUserScreen(context context.Context) SelectUserScreen {
	return SelectUserScreen{
		context:  context,
	}
}

// Init
func (screen SelectUserScreen) Init() tea.Cmd {
	return nil
}


// Update
func (screen SelectUserScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	// var cmd tea.Cmd

	// List of cmds to send
	cmds := make([]tea.Cmd, 0)

	// Handle key presses
	switch msg.(type) {
	case tea.KeyMsg:
		switch {

		}
	}

    // Update listview
	// screen.listview, cmd = screen.listview.Update(msg)
	// cmds = append(cmds, cmd)

    // Return
	return screen, tea.Batch(cmds...)
}

// View
func (screen SelectUserScreen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(selectUserAscii), "",
		tlockstyles.Styles.SubText.Render("Select a user to login as"), "",
	)
}
