package auth

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-internal/utils"
	"github.com/eklairs/tlock/tlock/models/dashboard"

	tlockcore "github.com/eklairs/tlock/tlock-core"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	tlockstyles "github.com/eklairs/tlock/tlock/styles"
)

// Select user ascii art
var selectUserAscii = `
█   █▀█ █▀▀ █ █▄ █
█▄▄ █▄█ █▄█ █ █ ▀█`

// Sudo ascii art
var sudoAscii = `
█▀ █ █ █▀▄ █▀█
▄█ █▄█ █▄▀ █▄█`

// select user list item
type selectUserListItem tlockcore.User

func (item selectUserListItem) FilterValue() string {
	return tlockcore.User(item).S()
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
	fmt.Fprint(w, renderer(65, string(item), "›"))
}

// Select user key map
type selectUserKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Enter   key.Binding
	New     key.Binding
	Options key.Binding
}

// ShortHelp()
func (k selectUserKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.New, k.Options, k.Enter}
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
	Options: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "user options"),
	),
}

// Select user
type SelectUserScreen struct {
	// Context
	context *context.Context

	// List view
	listview list.Model
}

// New instance of select user
func InitializeSelectUserScreen(context *context.Context) SelectUserScreen {
	// Renderable list of users
	usersList := utils.Map(context.Core.Users, func(user tlockcore.User) list.Item { return selectUserListItem(user) })

	// Return instance
	return SelectUserScreen{
		context:  context,
		listview: components.ListViewSimple(usersList, selectUserDelegate{}, 65, min(12, len(usersList)*3)),
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

	// Handle key presses
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		// New user
		case key.Matches(msgType, selectUserKeys.New):
			cmds = append(cmds, manager.PushScreen(InitializeCreateUserScreen(screen.context)))

        // User options
		case key.Matches(msgType, selectUserKeys.Options):
            // Try to unlock vault
            focused, vault := screen.tryUnlock();

            // If the vault is protected, ask for password
            if vault == nil {
                // Screen to go to
                next := InitializeEnterPassScreenCustomOpts(screen.context, tlockcore.User(focused), InitializeUserOptionsScreen, sudoAscii, "Enter password for %s to see user options")

                // Push
                cmds = append(cmds, manager.PushScreen(next))
            } else {
                cmds = append(cmds, manager.PushScreen(InitializeUserOptionsScreen(focused.S(), vault, screen.context)))
			}

		case key.Matches(msgType, selectUserKeys.Enter):
            // Try to unlock vault
            focused, vault := screen.tryUnlock();

            if vault == nil {
                // It is encrypted with a password, require password
                cmds = append(cmds, manager.PushScreen(InitializeEnterPassScreen(screen.context, focused, dashboard.InitializeDashboardScreen)))
            } else {
                // YAY!
                cmds = append(cmds, manager.PushScreen(dashboard.InitializeDashboardScreen(focused.S(), vault, screen.context)))
            }
		}

    case modelmanager.ScreenRefocusedMsg:
        // Update items
        screen.listview.SetItems(utils.Map(screen.context.Core.Users, func(user tlockcore.User) list.Item { return selectUserListItem(user) }))
	}

	// Update listview
	screen.listview, cmd = screen.listview.Update(msg)
	cmds = append(cmds, cmd)

	// Return
	return screen, tea.Batch(cmds...)
}

// View
func (screen SelectUserScreen) View() string {
    // Set height
    screen.listview.SetHeight(min(12, len(screen.context.Core.Users)*3))

	// List of items to render
	items := []string{
		tlockstyles.Title(selectUserAscii), "",
		tlockstyles.Dimmed("Select a user to login as"), "",
		screen.listview.View(), "",
	}

	// Add paginator
	if screen.listview.Paginator.TotalPages > 1 {
		items = append(items, components.Paginator(screen.listview), "")
	}

	// Add help
	items = append(items, tlockstyles.HelpView(selectUserKeys))

	// Return
	return lipgloss.JoinVertical(
		lipgloss.Center,
		items...,
	)
}

// Tries to unlock the vault for the focused user
func (screen SelectUserScreen) tryUnlock() (tlockcore.User, *tlockvault.Vault) {
    // Get focused
    focused := tlockcore.User(screen.listview.SelectedItem().(selectUserListItem))

    // Try to decrypt user with empty password
    vault, _ := tlockvault.Load(focused.Vault(), "")

    // Return
    return focused, vault
}
