package tokens

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/messages"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

// Folder list item
type moveTokenListItem tlockvault.Folder

func (item moveTokenListItem) FilterValue() string {
	return item.Name
}

// Folder list view delegate
type moveTokenDelegate struct{}

// Height
func (delegate moveTokenDelegate) Height() int {
	return 3
}

// Spacing
func (delegate moveTokenDelegate) Spacing() int {
	return 0
}

// Update
func (d moveTokenDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

// Render
func (d moveTokenDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(moveTokenListItem)

	if !ok {
		return
	}

	// Decide the renderer based on focused index
	renderer := components.ListItemInactive

	if index == m.Index() {
		renderer = components.ListItemActive
	}

	// Render
	fmt.Fprint(w, renderer(65, item.Name, "›"))
}

var moveTokenAscii = `
█▀▄▀█ █▀█ █ █ █▀▀
█ ▀ █ █▄█ ▀▄▀ ██▄`

// Move token key map
type moveTokenKeyMap struct {
	GoBack key.Binding
	Move   key.Binding
}

// ShortHelp()
func (k moveTokenKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.GoBack, k.Move}
}

// FullHelp()
func (k moveTokenKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.GoBack},
		{k.Move},
	}
}

// Keys
var moveTokenKeys = moveTokenKeyMap{
	Move: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "move"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

type MoveTokenScreen struct {
	// Vault
	vault *tlockvault.Vault

	// Token to move
	token tlockvault.Token

	// Folder
	folder tlockvault.Folder

	// Listview
	listview list.Model
}

// Initialize root model
func InitializeMoveTokenScreen(vault *tlockvault.Vault, folder tlockvault.Folder, token tlockvault.Token) MoveTokenScreen {
	items := make([]list.Item, len(vault.Folders))

	for index, folder := range vault.Folders {
		items[index] = moveTokenListItem(folder)
	}

	return MoveTokenScreen{
		vault:    vault,
		token:    token,
		folder:   folder,
		listview: components.ListViewSimple(items, moveTokenDelegate{}, 65, 15),
	}
}

// Init
func (screen MoveTokenScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen MoveTokenScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, moveTokenKeys.GoBack):
			manager.PopScreen()

		case key.Matches(msgType, moveTokenKeys.Move):
			focusedFolder := screen.listview.Items()[screen.listview.Index()].(moveTokenListItem)

			// Move token
			screen.vault.MoveToken(screen.token.ID, screen.folder.ID, focusedFolder.ID)

			// Require refresh of folders and tokens list
			cmds = append(
				cmds,
				func() tea.Msg { return tlockmessages.RefreshFoldersMsg{} },
				func() tea.Msg { return tlockmessages.RefreshTokensMsg{} },
			)

			// Pop
			manager.PopScreen()
		}
	}

	screen.listview, _ = screen.listview.Update(msg)

	return screen, tea.Batch(cmds...)
}

// View
func (screen MoveTokenScreen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(moveTokenAscii), "",
		tlockstyles.Styles.SubText.Render("Select the folder to move the token to"), "",
		screen.listview.View(), "",
		tlockstyles.Help.View(moveTokenKeys),
	)
}
