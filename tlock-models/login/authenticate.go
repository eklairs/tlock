package login

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	"github.com/eklairs/tlock/tlock-internal/buildhelp"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-models/dashboard"
	tlockvault "github.com/eklairs/tlock/tlock-vault"

	tea "github.com/charmbracelet/bubbletea"
	tlockcore "github.com/eklairs/tlock/tlock-core"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

var ENTER_PASS_SIZE = 65

// Enter pass key map
type enterPassKeyMap struct {
	Login key.Binding
	Back  key.Binding
}

// ShortHelp()
func (k enterPassKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Login, k.Back}
}

// FullHelp()
func (k enterPassKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Login},
		{k.Back},
	}
}

// Enter pass keys
var enterPassKeys = enterPassKeyMap{
	Login: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "login"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

// Enter pass Model
type EnterPassModel struct {
	// Context
	context context.Context

	// Common styles
	styles tlockstyles.Styles

	// Password input
	passInput textinput.Model

	// Help
	help help.Model

	// User spec
	userSpec tlockcore.UserSpec

	// Any error message
	errorMessage bool
}

// Initialize root model
func InitializeEnterPassModel(context context.Context, userSpec tlockcore.UserSpec) EnterPassModel {
	// Initialize styles
	styles := tlockstyles.InitializeStyle(ENTER_PASS_SIZE, context.Theme)

	// Password input
	passwordInput := tlockstyles.InitializeInputBox(styles, "Your password goes here...")
	passwordInput.Focus()

	// Help menu
	help := buildhelp.BuildHelp(styles)

	return EnterPassModel{
		context:   context,
		help:      help,
		userSpec:  userSpec,
		passInput: passwordInput,
		styles:    styles,
	}
}

// Init
func (model EnterPassModel) Init() tea.Cmd {
	return nil
}

// Update
func (model EnterPassModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	var cmd tea.Cmd

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, enterPassKeys.Back):
			manager.PopScreen()
		case key.Matches(msgType, enterPassKeys.Login):
			vault, err := tlockvault.Load(model.userSpec.Vault, model.passInput.Value())

			if err != nil {
				model.errorMessage = true
			} else {
				cmd = manager.PushScreen(dashboard.InitializeDashboardModel(*vault, model.context))
			}
		}
	}

	// Update input box
	model.passInput, _ = model.passInput.Update(msg)

	return model, cmd
}

// View
func (m EnterPassModel) View() string {
	// List of items
	items := []string{
		m.styles.Dimmed.Copy().AlignHorizontal(lipgloss.Center).Render(fmt.Sprintf("Login in as %s", m.userSpec.Username)), "",
		m.styles.Title.Render("Password"),
		m.styles.Dimmed.Render("Enter the super secret password"),
		m.styles.Input.Render(m.passInput.View()),
	}

	// Add error message if any
	if m.errorMessage {
		items = append(items, "", m.styles.Error.Render("Invalid password, please check if it is correct"))
	}

	// Rest
	items = append(items, "", m.styles.Center.Render(m.help.View(enterPassKeys)))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		items...,
	)
}
