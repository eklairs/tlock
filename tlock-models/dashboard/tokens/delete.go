package tokens

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	"github.com/pquerna/otp"
)

var deleteTokenAsciiArt = `
█▀▄ █▀▀ █   █▀▀ ▀█▀ █▀▀
█▄▀ ██▄ █▄▄ ██▄  █  ██▄`

type DeleteTokenMsg struct {
	FolderName string
    TokenURI string
}

// Delete token model key bindings
type deleteTokenKeyMap struct {
	Delete key.Binding
	GoBack key.Binding
}

// ShortHelp()
func (k deleteTokenKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.GoBack, k.Delete}
}

// FullHelp()
func (k deleteTokenKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.GoBack},
		{k.Delete},
	}
}

// Keys
var deleteTokenKeys = deleteTokenKeyMap{
	Delete: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "delete"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

// Delete token screen
type DeleteTokenScreen struct {
	// Token to delete
	token string

    // Folder
    folder string
}

// Initialize root model
func InitializeDeleteTokenScreen(folder, token string) DeleteTokenScreen {
	// Return
	return DeleteTokenScreen{
        folder: folder,
		token: token,
	}
}

// Init
func (screen DeleteTokenScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen DeleteTokenScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, deleteTokenKeys.GoBack):
			manager.PopScreen()
		case key.Matches(msgType, deleteTokenKeys.Delete):
			cmds = append(cmds, func() tea.Msg {
				return DeleteTokenMsg{
					FolderName: screen.folder,
                    TokenURI: screen.token,
				}
			})

			manager.PopScreen()
		}
	}

	return screen, tea.Batch(cmds...)
}

// View
func (screen DeleteTokenScreen) View() string {
    key, _ := otp.NewKeyFromURL(screen.token)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(deleteTokenAsciiArt), "",
		tlockstyles.Styles.SubText.Render("Permanently delete token"), "",
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			tlockstyles.Styles.SubText.Render("Are you sure you want to "),
			tlockstyles.Styles.Error.Render("× DELETE "),
			tlockstyles.Styles.Title.Render(key.AccountName()),
			tlockstyles.Styles.SubText.Render(" token ?"),
		), "",
		tlockstyles.Help.View(deleteTokenKeys),
	)
}
