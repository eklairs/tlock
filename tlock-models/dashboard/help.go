package dashboard

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	"golang.org/x/term"
)

var helpAsciiArt = `
█ █ █▀▀ █   █▀█
█▀█ ██▄ █▄▄ █▀▀`

type HelpKeyBindingSpec struct {
	// Key
	Key string

	// Description
	Desc string
}

type helpKeyBindings struct {
	// Keybindings for folders
	Folders []HelpKeyBindingSpec

	// Tokens
	Tokens []HelpKeyBindingSpec

	// Others
	Others []HelpKeyBindingSpec
}

// Builds the help menu for the given set of key bindings and the tile
func BuildHelpItem(title string, keys []HelpKeyBindingSpec) string {
	items := make([]string, 0)

	// Add title
	items = append(items, tlockstyles.Styles.Title.Render(title), "")

	// Add keys
	for _, key := range keys {
		ui := lipgloss.JoinHorizontal(
			lipgloss.Center,
			tlockstyles.Styles.SubText.Render(key.Desc),
			strings.Repeat(" ", 65-len(key.Desc)-len(key.Key)),
			tlockstyles.Styles.Title.Render(key.Key),
		)

		items = append(items, ui, "")
	}

	// Return
	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

func BuildHelpMenu(context *context.Context) string {
	m := func(keys []string) string {
		return strings.Join(keys, "/")
	}

	var helpKeys = helpKeyBindings{
		Folders: []HelpKeyBindingSpec{
			{
				Key:  m(context.Config.Folder.Add.Keys()),
				Desc: "Add a new folder",
			},
			{
				Key:  m(context.Config.Folder.Edit.Keys()),
				Desc: "Edit the current focused folder",
			},
			{
				Key:  m(context.Config.Folder.Next.Keys()),
				Desc: "Switch to next folder",
			},
			{
				Key:  m(context.Config.Folder.Previous.Keys()),
				Desc: "Switch to previous folder",
			},
			{
				Key:  m(context.Config.Folder.MoveUp.Keys()),
				Desc: "Move the focused folder up",
			},
			{
				Key:  m(context.Config.Folder.MoveDown.Keys()),
				Desc: "Move the focused folder down",
			},
			{
				Key:  m(context.Config.Folder.Delete.Keys()),
				Desc: "Delete the current focused folder",
			},
		},
		Tokens: []HelpKeyBindingSpec{
			{
				Key:  m(context.Config.Tokens.Add.Keys()),
				Desc: "Add a new token in the current focused folder",
			},
			{
				Key:  m(context.Config.Tokens.AddScreen.Keys()),
				Desc: "Add a new token from the screen",
			},
			{
				Key:  m(context.Config.Tokens.Edit.Keys()),
				Desc: "Edit the current focused token",
			},
			{
				Key:  m(context.Config.Tokens.Move.Keys()),
				Desc: "Move the current focused token to another folder",
			},
			{
				Key:  m(context.Config.Tokens.NextHOTP.Keys()),
				Desc: "Generates the token for the next counter [only of HOTP tokens]",
			},
			{
				Key:  m(context.Config.Tokens.Copy.Keys()),
				Desc: "Copy the current code for the focused token",
			},
			{
				Key:  m(context.Config.Tokens.Next.Keys()),
				Desc: "Move focus to the next token",
			},
			{
				Key:  m(context.Config.Tokens.Previous.Keys()),
				Desc: "Move focus to the previous token",
			},
			{
				Key:  m(context.Config.Tokens.MoveUp.Keys()),
				Desc: "Move the focused token up",
			},
			{
				Key:  m(context.Config.Tokens.MoveDown.Keys()),
				Desc: "Move the focused token down",
			},
			{
				Key:  m(context.Config.Tokens.Delete.Keys()),
				Desc: "Delete the current focused token",
			},
		},
		Others: []HelpKeyBindingSpec{
			{
				Key:  "?",
				Desc: "Show this help window",
			},
			{
				Key:  "ctrl+t",
				Desc: "Change theme",
			},
			{
				Key:  "ctrl+c / ctrl+q",
				Desc: "Exit the application",
			},
		},
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(helpAsciiArt), "",
		tlockstyles.Styles.SubText.Render("Keybindings to move around the app"), "",
		BuildHelpItem("Folders", helpKeys.Folders),
		BuildHelpItem("Tokens", helpKeys.Tokens),
		BuildHelpItem("Others", helpKeys.Others),
	)
}

// Help screen
type HelpScreen struct {
	viewport viewport.Model
}

// Initializes a new instance of the help screen
func InitializeHelpScreen(context *context.Context) HelpScreen {
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))

	viewport := viewport.New(65, height)
	viewport.SetContent(BuildHelpMenu(context))

	return HelpScreen{
		viewport: viewport,
	}
}

// Init
func (screen HelpScreen) Init() tea.Cmd {
	return screen.viewport.Init()
}

// Update
func (screen HelpScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "esc":
			manager.PopScreen()
		}
	case tea.WindowSizeMsg:
		screen.viewport.Height = msgType.Height
	}

	// Update viewport
	screen.viewport, _ = screen.viewport.Update(msg)

	return screen, nil
}

// View
func (screen HelpScreen) View() string {
	return screen.viewport.View()
}
