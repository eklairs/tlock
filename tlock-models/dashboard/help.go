package dashboard

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	"golang.org/x/term"
)

type HelpKeyBinding struct {
	// Key
	Key string

	// Description
	Desc string
}

var HELP_WIDTH = 65

var helpAsciiArt = `
█ █ █▀▀ █   █▀█
█▀█ ██▄ █▄▄ █▀▀`

var helpKeys = map[string][]HelpKeyBinding{
	"Folders": {
		{
			Key:  "A",
			Desc: "Add a new folder",
		},
		{
			Key:  "E",
			Desc: "Edit the current focused folder",
		},
		{
			Key:  "Shift + Up",
			Desc: "Move the focused folder up",
		},
		{
			Key:  "Shift + Down",
			Desc: "Move the focused folder down",
		},
		{
			Key:  "D",
			Desc: "Delete the current focused folder",
		},
	},
	"Tokens": {
		{
			Key:  "a",
			Desc: "Add a new token in the current focused folder",
		},
		{
			Key:  "e",
			Desc: "Edit the current focused token",
		},
		{
			Key:  "m",
			Desc: "Move the current focused token to another folder",
		},
		{
			Key:  "c",
			Desc: "Copy the current code for the focused token",
		},
		{
			Key:  "Ctrl + Up",
			Desc: "Move the focused token up",
		},
		{
			Key:  "Ctrl + Down",
			Desc: "Move the focused token down",
		},
		{
			Key:  "Ctrl + /",
			Desc: "Search for a token inside the focused folder",
		},
		{
			Key:  "D",
			Desc: "Delete the current focused tokens",
		},
	},
	"Others": {
		{
			Key:  "?",
			Desc: "Show this help message",
		},
		{
			Key:  "Ctrl + T",
			Desc: "Change theme",
		},
		{
			Key:  "Ctrl + C / Ctrl + Q",
			Desc: "Exit the application",
		},
	},
}

func BuildHelpMenu(styles tlockstyles.Styles) string {
	items := make([]string, 0)

	// Title
	items = append(items, styles.Center.Render(styles.Title.Copy().UnsetWidth().Render(helpAsciiArt)), "")

	// Some description
	items = append(items, styles.Center.Render(styles.Dimmed.Render("Keybindings to move around the app")), "", "")

	for sectionName, keys := range helpKeys {
		items = append(items, "", styles.Title.Render(sectionName), "")

		for _, key := range keys {
			ui := lipgloss.JoinHorizontal(
				lipgloss.Center,
				styles.Dimmed.Copy().UnsetWidth().Render(key.Desc),
				strings.Repeat(" ", HELP_WIDTH-len(key.Desc)-len(key.Key)),
				styles.Title.Copy().UnsetWidth().Render(key.Key),
			)

			items = append(items, ui, "")
		}
	}

	return lipgloss.JoinVertical(lipgloss.Center, items...)
}

// HelpModel struct
type HelpModel struct {
	// Styles
	Styles tlockstyles.Styles

	// Viewport
	Viewport viewport.Model
}

// Initializes a new instance of help model
func InitializeHelpModel(context context.Context) HelpModel {
	_, height, _ := term.GetSize(0)
	styles := tlockstyles.InitializeStyle(HELP_WIDTH, context.Theme)

	viewport := viewport.New(HELP_WIDTH, height)
	viewport.SetContent(BuildHelpMenu(styles))

	return HelpModel{
		Styles:   styles,
		Viewport: viewport,
	}
}

// Init()
func (model HelpModel) Init() tea.Cmd {
	return nil
}

// Update()
func (model HelpModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "esc":
			manager.PopScreen()
		}
	}

	model.Viewport, _ = model.Viewport.Update(msg)

	return model, nil
}

// View()
func (model HelpModel) View() string {
	return model.Viewport.View()
}
