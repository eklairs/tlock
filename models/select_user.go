package models

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/internal/boundedinteger"

	tea "github.com/charmbracelet/bubbletea"
	tlockcore "github.com/eklairs/tlock/tlock-core"

	. "github.com/eklairs/tlock/internal/modelmanager"
)

var SELECT_USER_SIZE = 60

// Select user key map
type selectUserKeyMap struct {
    Up key.Binding
    Down key.Binding
    New key.Binding
}

// ShortHelp()
func (k selectUserKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.New}
}

// LongHelp()
func (k selectUserKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
        {k.Up},
        {k.Down},
        {k.New},
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
}

// Select user ascii
var selectUserAscii = `
 ____       _           _     _   _               
/ ___|  ___| | ___  ___| |_  | | | |___  ___ _ __ 
\___ \ / _ \ |/ _ \/ __| __| | | | / __|/ _ \ '__|
 ___) |  __/ |  __/ (__| |_  | |_| \__ \  __/ |   
|____/ \___|_|\___|\___|\__|  \___/|___/\___|_|   
`

// Select User Model
type SelectUserModel struct {
    // Styles
    styles Styles

    // Help
    help help.Model

    // Core
    core tlockcore.TLockCore

    // Focused user index
    focused_index boundedinteger.BoundedInteger
}

// Initializes a new instance of the SelectUserModel
func InitializeSelectUserModel(core tlockcore.TLockCore) SelectUserModel {
    return SelectUserModel {
        core: core,
        help: help.New(),
        styles: InitializeStyles(SELECT_USER_SIZE),
        focused_index: boundedinteger.New(0, len(core.Users.Users)),
    }
}

// Init
func (m SelectUserModel) Init() tea.Cmd {
    return nil
}

// Update
func (m SelectUserModel) Update(msg tea.Msg, manager *ModelManager) (Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msgType, selectUserKeys.Down):
            m.focused_index.Increase()
        
        case key.Matches(msgType, selectUserKeys.Up):
            m.focused_index.Decrease()

        case key.Matches(msgType, selectUserKeys.New):
            manager.PushScreen(InitializeNewUserModel(m.core))
        }
    }

	return m, nil
}

// View
func (m SelectUserModel) View() string {
    // List of ui items
    items := []string {
        m.styles.dimmedCenter.Render("Select a user to continue as"), "",
    }

    // Render all the list users
    for index, user := range m.core.Users.Users {
        render_fn := m.styles.inactive.Render

        if index == m.focused_index.Value {
            render_fn = m.styles.active.Render
        }

        renderable := render_fn(
            lipgloss.JoinHorizontal(
                lipgloss.Center,
                user.Username,
                strings.Repeat(" ", SELECT_USER_SIZE - len(user.Username) - 1 - 6),
                "›",
            ),
        )

        items = append(items, renderable)
    }

    // Help
    items = append(items, "", m.styles.center.Render(m.help.View(selectUserKeys)))

    // Join them all!
    return lipgloss.JoinVertical(
        lipgloss.Left,
        items...
    )
}

