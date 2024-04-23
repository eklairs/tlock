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

type EditNewFolderMsg struct {
	OldName string
	NewName string
}

var EDIT_FOLDER_SIZE = 65

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

// Edit folder Model
type EditFolderModel struct {
	// Styles
	styles tlockstyles.Styles

	// Folder name input
	name textinput.Model

	// Help
	help help.Model

	// Original name
	oldName string
}

// Initialize edit folder model
func InitializeEditFolderModel(oldName string, context context.Context) EditFolderModel {
	// Initialize style
	styles := tlockstyles.InitializeStyle(ADD_FOLDER_SIZE, context.Theme)

	// Initialize input box
	name := tlockstyles.InitializeInputBox(styles, "Your new folder name goes here...")
	name.SetValue(oldName)
	name.Focus()

	// Help menu
	help := buildhelp.BuildHelp(styles)

	return EditFolderModel{
		name:    name,
		help:    help,
		styles:  styles,
		oldName: oldName,
	}
}

// Init
func (model EditFolderModel) Init() tea.Cmd {
	return nil
}

// Update
func (model EditFolderModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, editFolderKeys.GoBack):
			manager.PopScreen()

		case key.Matches(msgType, editFolderKeys.Enter):
			cmds = append(cmds, func() tea.Msg {
				return EditNewFolderMsg{
					OldName: model.oldName,
					NewName: model.name.Value(),
				}
			})

			manager.PopScreen()
		}
	}

	model.name, _ = model.name.Update(msg)

	return model, tea.Batch(cmds...)
}

// View
func (m EditFolderModel) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.Center.Render(m.styles.Title.Render(editFolderAscii)), "",
		m.styles.Center.Render(m.styles.Dimmed.Render("Rename a folder to a new name")), "",
		m.styles.Title.Render("Name"),
		m.styles.Dimmed.Render("Choose the new name for your folder"), "",
		m.styles.Input.Render(m.name.View()), "",
		m.styles.Center.Render(m.help.View(editFolderKeys)),
	)
}
