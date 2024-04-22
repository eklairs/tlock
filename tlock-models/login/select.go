package login

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	tlockstyles "github.com/eklairs/tlock/tlock-styles"

	"github.com/eklairs/tlock/tlock-internal/boundedinteger"
	"github.com/eklairs/tlock/tlock-internal/buildhelp"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
)

var SELECT_USER_WIDTH = 65

var selectUserAsciiArt = `
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

// Select user screen
type SelectUserModel struct {
	// Context
	context context.Context

	// Styles
	styles tlockstyles.Styles

	// Focused index
	focused_index boundedinteger.BoundedInteger

	// Help
	help help.Model
}

// Initializes a new instance of the select user model
func InitializeSelectUserModel(context context.Context) SelectUserModel {
	// Initialize styles
	styles := tlockstyles.InitializeStyle(SELECT_USER_WIDTH, context.Theme)

	// Return
	return SelectUserModel{
		context:       context,
		help:          buildhelp.BuildHelp(styles),
		focused_index: boundedinteger.New(0, len(context.Core.Users)),
		styles:        styles,
	}
}

// Init
func (model SelectUserModel) Init() tea.Cmd {
	return nil
}

// Update
func (model SelectUserModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, selectUserKeys.Up):
			model.focused_index.Decrease()
		case key.Matches(msgType, selectUserKeys.Down):
			model.focused_index.Increase()
		case key.Matches(msgType, selectUserKeys.New):
			manager.PushScreen(InitializeCreateUserModel(model.context))
		case key.Matches(msgType, selectUserKeys.Enter):
			manager.PushScreen(InitializeEnterPassModel(model.context, model.context.Core.Users[model.focused_index.Value]))
		}
	}

	return model, nil
}

// View
func (model SelectUserModel) View() string {
	// Inner width of each item
	// -6 is for the padding on both the ends of the horizontal axis
	listItemWidth := SELECT_USER_WIDTH - 6

	items := []string{
		model.styles.Title.Copy().UnsetWidth().Render(selectUserAsciiArt), "",
		model.styles.Dimmed.Copy().UnsetWidth().Render("Select a user to continue as"), "",
	}

	// Add all the user's name
	for index, user := range model.context.Core.Users {
		render_fn := model.styles.InactiveListItem.Render

		if index == model.focused_index.Value {
			render_fn = model.styles.ActiveItem.Render
		}

		// Join
		ui := lipgloss.JoinHorizontal(
			lipgloss.Center,
			user.Username,
			strings.Repeat(" ", listItemWidth-1-len(user.Username)),
			"›",
		)

		items = append(items, render_fn(ui))
	}

	// Add the help menu
	items = append(items, "", model.styles.Center.Render(model.help.View(selectUserKeys)))

	return lipgloss.JoinVertical(
		lipgloss.Center,
		items...,
	)
}
