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

var ERROR_EMPTY_FOLDER_NAME = "Folder name cannot be empty"

// Message that represents adding a new folder
type AddNewFolderMsg struct {
	FolderName string
}

var addFolderAscii = `
▄▀█ █▀▄ █▀▄
█▀█ █▄▀ █▄▀`

// Add folder key map
type addFolderKeyMap struct {
	GoBack key.Binding
	Enter  key.Binding
}

// ShortHelp()
func (k addFolderKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.GoBack, k.Enter}
}

// FullHelp()
func (k addFolderKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.GoBack},
		{k.Enter},
	}
}

// Keys
var addFolderKeys = addFolderKeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "create folder"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

// Add folder screen
type AddFolderScreen struct {
	// Folder name input
	name textinput.Model

	// Help
	help help.Model

	// Error
	errorMessage *string
}

// Initialize add folder scree
func InitializeAddFolderScreen() AddFolderScreen {
	// Initialize input box
	name := components.InitializeInputBox("Your folder name goes here...")
	name.Focus()

	// Return
	return AddFolderScreen{
		name: name,
		help: components.BuildHelp(),
	}
}

// Init
func (screen AddFolderScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen AddFolderScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, addFolderKeys.GoBack):
			manager.PopScreen()

		case key.Matches(msgType, addFolderKeys.Enter):
			if screen.name.Value() == "" {
				screen.errorMessage = &ERROR_EMPTY_FOLDER_NAME
			} else {
				cmds = append(cmds, func() tea.Msg {
					return AddNewFolderMsg{
						FolderName: screen.name.Value(),
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
func (screen AddFolderScreen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(addFolderAscii), "",
		tlockstyles.Styles.SubText.Render("Add a new folder"), "",
		components.InputGroup("Name", "Choose a name for your folder, like Socials!", screen.errorMessage, screen.name),
		screen.help.View(addFolderKeys),
	)
}
