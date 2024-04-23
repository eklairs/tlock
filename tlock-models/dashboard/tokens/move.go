package tokens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/boundedinteger"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

var moveTokenAscii = `
█▀▄▀█ █▀█ █ █ █▀▀
█ ▀ █ █▄█ ▀▄▀ ██▄`

type MoveTokenModel struct {
    // Styles
    styles tlockstyles.Styles

    // Vault
    vault *tlockvault.TLockVault

    // Token to move
    token string

    // Folder
    folder string

    // Focused index
    focused_index boundedinteger.BoundedInteger
}

// Initialize root model
func InitializeMoveTokenModel(vault *tlockvault.TLockVault, context context.Context, folder, token string) MoveTokenModel {
    return MoveTokenModel {
        vault: vault,
        token: token,
        folder: folder,
        styles: tlockstyles.InitializeStyle(65, context.Theme),
        focused_index: boundedinteger.New(0, len(vault.Data.Folders)),
    }
}


// Init
func (m MoveTokenModel) Init() tea.Cmd {
    return nil
}

// Update
func (m MoveTokenModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch msgType.String() {
        case tea.KeyDown.String(), "j":
            m.focused_index.Increase()
        case tea.KeyUp.String(), "k":
            m.focused_index.Decrease()
        case "esc":
            manager.PopScreen()
        case "enter":
            m.vault.MoveURI(m.folder, m.token, m.vault.Data.Folders[m.focused_index.Value].Name)

            manager.PopScreen()
        }
    }

	return m, nil
}

// View
func (m MoveTokenModel) View() string {
    user_items := []string {
        m.styles.Center.Render(m.styles.Title.Render(moveTokenAscii)), "",
        m.styles.Center.Render(m.styles.Dimmed.Render("Select the folder to move the token to")), "",
    }

    for index, folder := range m.vault.Data.Folders {
        render_fn := m.styles.InactiveListItem.Render

        if index == m.focused_index.Value {
            render_fn = m.styles.ActiveItem.Render
        }

        renderable := render_fn(
            lipgloss.JoinHorizontal(
                lipgloss.Center,
                folder.Name,
                strings.Repeat(" ", 65 - len(folder.Name) - 8 - 6),
                fmt.Sprintf("%d tokens", len(folder.Uris)),
            ),
        )

        user_items = append(user_items, renderable)
    }

    return lipgloss.JoinVertical(
        lipgloss.Left,
        user_items...
    )
}

