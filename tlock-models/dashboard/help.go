package dashboard

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

var helpKeys = helpKeyBindings{
	Folders: []HelpKeyBindingSpec{
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
	Tokens: []HelpKeyBindingSpec{
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
			Key:  "D",
			Desc: "Delete the current focused tokens",
		},
	},
	Others: []HelpKeyBindingSpec{
		{
			Key:  "?",
			Desc: "Show this help window",
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

func BuildHelpMenu() string {
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
func InitializeHelpScreen() HelpScreen {
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))

	viewport := viewport.New(65, height)
	viewport.SetContent(BuildHelpMenu())

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
	}

	// Update viewport
	screen.viewport.Update(msg)

	return screen, nil
}

// View
func (screen HelpScreen) View() string {
	return screen.viewport.View()
}
