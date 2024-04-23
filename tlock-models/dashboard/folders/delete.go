package folders

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/buildhelp"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

var DELETE_FOLDER_SIZE = 65

var deleteFolderAsciiArt = `
█▀▄ █▀▀ █   █▀▀ ▀█▀ █▀▀
█▄▀ ██▄ █▄▄ ██▄  █  ██▄`

type DeleteFolderMsg struct {
	FolderName string
}

// Create user model key bindings
type deleteFolderKeyMap struct {
	Delete key.Binding
	GoBack key.Binding
}

// ShortHelp()
func (k deleteFolderKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.GoBack, k.Delete}
}

// FullHelp()
func (k deleteFolderKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.GoBack},
		{k.Delete},
	}
}

// Keys
var deleteFolderKeys = deleteFolderKeyMap{
	Delete: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "delete"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

// Delete folder Model
type DeleteFolderModel struct {
	// Styles
	styles tlockstyles.Styles

	// Name of the folder to be deleted
	folder string

	// Help
	help help.Model
}

// Initialize root model
func InitializeDeleteFolderModel(folder string, context context.Context) DeleteFolderModel {
	// Initialize styles
	styles := tlockstyles.InitializeStyle(DELETE_FOLDER_SIZE, context.Theme)

	// Initialize help
	help := buildhelp.BuildHelp(styles)

	return DeleteFolderModel{
		styles: styles,
		folder: folder,
		help:   help,
	}
}

// Init
func (model DeleteFolderModel) Init() tea.Cmd {
	return nil
}

// Update
func (model DeleteFolderModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, deleteFolderKeys.GoBack):
			manager.PopScreen()
		case key.Matches(msgType, deleteFolderKeys.Delete):
			cmds = append(cmds, func() tea.Msg {
				return DeleteFolderMsg{
					FolderName: model.folder,
				}
			})

			manager.PopScreen()
		}
	}

	return model, tea.Batch(cmds...)
}

// View
func (model DeleteFolderModel) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		model.styles.Center.Render(model.styles.Title.Render(deleteFolderAsciiArt)), "",
		model.styles.Center.Render(model.styles.Dimmed.Render("Permanently delete tokens folder")), "",
		model.styles.Center.Render(lipgloss.JoinHorizontal(
			lipgloss.Center,
			model.styles.Base.Copy().UnsetWidth().Render("Are you sure you want to "),
			model.styles.Error.Copy().UnsetWidth().Render("DELETE "),
			model.styles.Base.Copy().UnsetWidth().Render("folder "),
			model.styles.Title.Copy().UnsetWidth().Render(model.folder),
			model.styles.Base.Copy().UnsetWidth().Render("?"),
		)),
		model.styles.Center.Render(model.styles.Base.Render("This action is irreversible!")), "",
		model.styles.Center.Render(model.help.View(deleteFolderKeys)),
	)
}
