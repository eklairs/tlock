package folders

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"

	tea "github.com/charmbracelet/bubbletea"
)

type EditFolderMsg struct {
	OldName string
	NewName string
}

var editFolderAscii = `
█▀▀ █▀▄ █ ▀█▀
██▄ █▄▀ █  █`

// Edit folder key map
type editFolderKeyMap struct {
	GoBack key.Binding
	Enter  key.Binding
}

// ShortHelp()
func (k editFolderKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.GoBack, k.Enter}
}

// FullHelp()
func (k editFolderKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.GoBack},
		{k.Enter},
	}
}

// Keys
var editFolderKeys = editFolderKeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "edit folder"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

// Edit folder screen
type EditFolderScreen struct {
	// Folder name input
	name textinput.Model

	// Help
	help help.Model

	// Original name
	oldName string

	// Error message
	errorMessage *string
}

// Initialize edit folder screen
func InitializeEditFolderScreen(oldName string) EditFolderScreen {
	// Initialize input box
	name := components.InitializeInputBox("Your folder name goes here...")
	name.SetValue(oldName)
	name.Focus()

	// Return
	return EditFolderScreen{
		name:    name,
		help:    components.BuildHelp(),
		oldName: oldName,
	}
}

// Init
func (screen EditFolderScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen EditFolderScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, editFolderKeys.GoBack):
			manager.PopScreen()

		case key.Matches(msgType, editFolderKeys.Enter):
			if screen.name.Value() == "" {
				screen.errorMessage = &ERROR_EMPTY_FOLDER_NAME
			} else {
				cmds = append(cmds, func() tea.Msg {
					return EditFolderMsg{
						NewName: screen.name.Value(),
						OldName: screen.oldName,
					}
				})

				manager.PopScreen()
			}
		}
	}

	screen.name, _ = screen.name.Update(msg)

	return screen, tea.Batch(cmds...)
}

// View
func (screen EditFolderScreen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(editFolderAscii), "",
		tlockstyles.Styles.SubText.Render("Rename the folder to a new name"), "",
		components.InputGroup("Name", "Choose the new name for your folder", screen.errorMessage, screen.name),
		screen.help.View(editFolderKeys),
	)
}
