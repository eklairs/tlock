package login

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/buildhelp"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-models/dashboard"

	tea "github.com/charmbracelet/bubbletea"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

var CREATE_USER_WIDTH = 65

// Create user key map
type createUserKeyMap struct {
	Tab    key.Binding
	Create key.Binding
	GoBack key.Binding
}

// ShortHelp()
func (k createUserKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Create, k.GoBack}
}

// FullHelp()
func (k createUserKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab},
		{k.Create},
		{k.GoBack},
	}
}

// Keys
var createUserKeys = createUserKeyMap{
	Tab: key.NewBinding(
		key.WithKeys("tab", "shift+tab"),
		key.WithHelp("tab/shift+tab", "switch input"),
	),
	Create: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "create"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

// Create user screen
type CreateUserModel struct {
	// Context
	context context.Context

	// Styles
	styles tlockstyles.Styles

	// Help
	help help.Model

	// Username input
	usernameInput textinput.Model

	// Password input
	passwordInput textinput.Model
}

// Initializes a new instance of the create user model
func InitializeCreateUserModel(context context.Context) CreateUserModel {
	// Initialize styles
	styles := tlockstyles.InitializeStyle(SELECT_USER_WIDTH, context.Theme)

	// Initialize help menu
	help := buildhelp.BuildHelp(styles)

	// Input box for username
	usernameInput := tlockstyles.InitializeInputBox(styles, "Your username goes here...")
	usernameInput.Focus()

	return CreateUserModel{
		context:       context,
		help:          help,
		styles:        styles,
		usernameInput: usernameInput,
		passwordInput: tlockstyles.InitializeInputBox(styles, "Your password goes here..."),
	}
}

// Init
func (model CreateUserModel) Init() tea.Cmd {
	return nil
}

// Update
func (model CreateUserModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, createUserKeys.GoBack):
			manager.PopScreen()

		case key.Matches(msgType, createUserKeys.Tab):
			if model.usernameInput.Focused() {
				model.usernameInput.Blur()
				model.passwordInput.Focus()
			} else {
				model.usernameInput.Focus()
				model.passwordInput.Blur()
			}

		case key.Matches(msgType, createUserKeys.Create):
			vault := model.context.Core.AddNewUser(model.usernameInput.Value(), model.passwordInput.Value())

			manager.PushScreen(dashboard.InitializeDashboardModel(vault, model.context))
		}
	}

	if model.usernameInput.Focused() {
		model.usernameInput, _ = model.usernameInput.Update(msg)
	}

	if model.passwordInput.Focused() {
		model.passwordInput, _ = model.passwordInput.Update(msg)
	}

	return model, nil
}

// View
func (model CreateUserModel) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		model.styles.Title.Render("Username"),
		model.styles.Dimmed.Render("Choose an awesome username, like Komaru!"),
		model.styles.Input.Render(model.usernameInput.View()), "",
		model.styles.Title.Render("Password"),
		model.styles.Dimmed.Render("Choose a super strong password, or keep it empty if you don't want any password"),
		model.styles.Input.Render(model.passwordInput.View()), "",
		model.styles.Center.Render(model.help.View(createUserKeys)),
	)
}
