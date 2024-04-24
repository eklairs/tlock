package auth

import (
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"

	tea "github.com/charmbracelet/bubbletea"
)

// Width of select user
var SELECT_USER_WIDTH = 65

// select user list item
type selectUserListItem string

func (item selectUserListItem) FilterValue() string {
	return string(item)
}

// Select user list view delegate
type selectUserDelegate struct{}

// Height
func (delegate selectUserDelegate) Height() int {
	return 3
}

// Spacing
func (delegate selectUserDelegate) Spacing() int {
	return 0
}

// Update
func (d selectUserDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

// Render
func (d selectUserDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {

}

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
	usersItem := make([]list.Item, len(context.Core.Users))

	// Iter
	for index, user := range context.Core.Users {
		usersItem[index] = selectUserListItem(user.Username)
	}

	// Build listview
	listview := components.ListViewSimple(usersItem, selectUserDelegate{}, SELECT_USER_WIDTH, 30)

	return SelectUserScreen{
		context:  context,
		listview: listview,
	}
}

// Init
func (screen SelectUserScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen SelectUserScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	var cmd tea.Cmd

	// List of cmds to send
	cmds := make([]tea.Cmd, 0)

	screen.listview, cmd = screen.listview.Update(msg)
	cmds = append(cmds, cmd)

	return screen, tea.Batch(cmds...)
}

// View
func (screen SelectUserScreen) View() string {
	return screen.listview.View()
}
