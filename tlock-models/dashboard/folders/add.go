package folders

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/buildhelp"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"

	tea "github.com/charmbracelet/bubbletea"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

type AddNewFolderMsg struct {
	FolderName string
}

var ADD_FOLDER_SIZE = 65

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

// Add folder Model
type AddFolderModel struct {
	// Styles
	styles tlockstyles.Styles

	// Folder name input
	name textinput.Model

	// Help
	help help.Model
}

// Initialize add folder model
func InitializeAddFolderModel(context context.Context) AddFolderModel {
	// Initialize theme
	styles := tlockstyles.InitializeStyle(ADD_FOLDER_SIZE, context.Theme)

	// Initialize input box
	name := tlockstyles.InitializeInputBox(styles, "Your folder name goes here...")
	name.Focus()

	// Help menu
	help := buildhelp.BuildHelp(styles)

	return AddFolderModel{
		name:   name,
		help:   help,
		styles: styles,
	}
}

// Init
func (model AddFolderModel) Init() tea.Cmd {
	return nil
}

// Update
func (model AddFolderModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, addFolderKeys.GoBack):
			manager.PopScreen()

		case key.Matches(msgType, addFolderKeys.Enter):
			cmds = append(cmds, func() tea.Msg {
				return AddNewFolderMsg{
					FolderName: model.name.Value(),
				}
			})

			manager.PopScreen()
		}
	}

	model.name, _ = model.name.Update(msg)

	return model, tea.Batch(cmds...)
}

// View
func (m AddFolderModel) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.Title.Render("Name"), // Folder name header
		m.styles.Dimmed.Render("Choose a name for your folder, like Socials!"), "", // Folder name description
		m.styles.Input.Render(m.name.View()), "",
		m.styles.Center.Render(m.help.View(addFolderKeys)),
	)
}
