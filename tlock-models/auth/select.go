package auth

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"

	tea "github.com/charmbracelet/bubbletea"
)

var selectUserAsciiArt = `
█   █▀█ █▀▀ █ █▄ █
█▄▄ █▄█ █▄█ █ █ ▀█`

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
	item, ok := listItem.(selectUserListItem)

	if !ok {
		return
	}

	// Decide the renderer based on focused index
	renderer := components.ListItemInactive
	if index == m.Index() {
		renderer = components.ListItemActive
	}

	// Render
	fmt.Fprint(w, renderer(SELECT_USER_WIDTH, string(item), "›"))
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

	// Help
	help help.Model
}

// New instance of select user
func InitializeSelectUserScreen(context context.Context) SelectUserScreen {
	usersItem := make([]list.Item, len(context.Core.Users))

	// Iter
	for index, user := range context.Core.Users {
		usersItem[index] = selectUserListItem(user.Username)
	}

	// Build listview
	listview := components.ListViewSimple(usersItem, selectUserDelegate{}, SELECT_USER_WIDTH, 13)

	return SelectUserScreen{
		context:  context,
		listview: listview,
		help:     components.BuildHelp(),
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
	return lipgloss.JoinVertical(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(selectUserAsciiArt), "",
		tlockstyles.Styles.SubText.Render("Select a user to login as"), "",
		screen.listview.View(),
		screen.help.View(selectUserKeys),
	)
}
