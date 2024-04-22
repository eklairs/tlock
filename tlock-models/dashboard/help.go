package dashboard

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

var HELP_WIDTH = 65

// HelpModel struct
type HelpModel struct {
	// Styles
	Styles tlockstyles.Styles
}

// Initializes a new instance of help model
func InitializeHelpModel(context context.Context) HelpModel {
	return HelpModel{
		Styles: tlockstyles.InitializeStyle(HELP_WIDTH, context.Theme),
	}
}

// Init()
func (model HelpModel) Init() tea.Cmd {
	return nil
}

// Update()
func (model HelpModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	return model, nil
}

// View()
func (model HelpModel) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		model.keyInfo("?", "Show this help window"),
		model.keyInfo("Ctrl+/", "Search for a token"),
		model.keyInfo("Ctrl+T", "Change theme"),
		model.keyInfo("A", "Add a new folder"),
		model.keyInfo("E", "Edit the current focused folder"),
		model.keyInfo("X", "Delete the current focused folder"),
		model.keyInfo("Shift+Up", "Move the folder up"),
		model.keyInfo("Shift+Down", "Move the folder down"),
		model.keyInfo("a", "Add a new token inside current focused folder"),
		model.keyInfo("e", "Edit the current focused token"),
		model.keyInfo("x", "Delete the current focused token"),
		model.keyInfo("u", "Generate code for the next counter [only for HOTP based tokens]"),
		model.keyInfo("Enter", "Copy the current focused token code to clipboard"),
		model.keyInfo("Ctrl+C", "Exit the application"),
		model.keyInfo("Ctrl+R", "Refresh the entire UI"),
	)
}

func (model HelpModel) keyInfo(key string, desc string) string {
	key_style := model.Styles.Title.Copy().Width(10)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		key_style.Render(key),
		strings.Repeat(" ", 4),
		model.Styles.Base.Render(desc),
	)
}
