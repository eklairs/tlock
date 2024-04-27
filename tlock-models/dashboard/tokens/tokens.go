package tokens

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tlockinternal "github.com/eklairs/tlock/tlock-internal"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/tlock-messages"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
	"golang.design/x/clipboard"
	"golang.org/x/term"
)

// Returns the width for the tokens listview for the given window wdith
func tokensWidth(width int) int {
	return width - int(math.Floor((1.0/5.0)*float64(width)))
}

// Returns the remaining time
func getRemainingTime(token tlockvault.Token) int {
	return int(token.Period - int(time.Now().Unix())%token.Period)
}

// Returns the current code
func getCurrentCode(tokenType tlockvault.TokenType, secret string, usageCounter int) string {
	var code string

	if tokenType == tlockvault.TokenTypeTOTP {
		code, _ = totp.GenerateCode(secret, time.Now())
	} else {
		code, _ = hotp.GenerateCode(secret, uint64(usageCounter))
	}

	return code
}

// Token list item
type tokensListItem struct {
	// Current code
	CurrentCode string

	// URI string
	Token tlockvault.Token

	// Time remaining before the otp is updated
	// Only in case of totp tokens
	time *int
}

func (item tokensListItem) FilterValue() string {
	return ""
}

// Refreshes the token
func (item *tokensListItem) Refresh() {
	// If the token is a totp, then update the time
	if item.Token.Type == tlockvault.TokenTypeTOTP {
		timeToRefresh := getRemainingTime(item.Token)
		item.time = &timeToRefresh
	}

	// Update current code
	item.CurrentCode = getCurrentCode(item.Token.Type, item.Token.Secret, item.Token.UsageCounter)
}

// Initializes a new instance of the tokens list item
func InitializeTokenListItem(token tlockvault.Token) tokensListItem {
	var ttr *int

	if token.Type == tlockvault.TokenTypeTOTP {
		timeToRefresh := getRemainingTime(token)
		ttr = &timeToRefresh
	}

	return tokensListItem{
		CurrentCode: getCurrentCode(token.Type, token.Secret, token.UsageCounter),
		Token:       token,
		time:        ttr,
	}
}

// Tokens key map
type tokenKeyMap struct {
	Manual key.Binding
	Screen key.Binding
}

// ShortHelp()
func (k tokenKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Manual, k.Screen}
}

// FullHelp()
func (k tokenKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Manual},
		{k.Screen},
	}
}

// Keys
var tokenKeys = tokenKeyMap{
	Manual: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add token"),
	),
	Screen: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "add from screen"),
	),
}

var EmptyAsciiArt = `
\    /\
 )  ( ')
(  /  )
 \(__)|
`

// Tokens list delegate
type tokensListDelegate struct{}

// Height
func (d tokensListDelegate) Height() int {
	return 4
}

// Spacing
func (d tokensListDelegate) Spacing() int {
	return 0
}

// Update
func (d tokensListDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

// Render
func (d tokensListDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item := listItem.(tokensListItem)

	// Decide renderer function
	render_fn := components.TokenItemInactive

	if index == m.Index() {
		render_fn = components.TokenItemActive
	}

	// Account name
	account := item.Token.Account

	if account == "" {
		account = "<no account name>"
	}

	// Issuer name
	issuer := item.Token.Issuer

	if issuer == "" {
		issuer = "<no issuer name>"
	}

	// Suffix (current code)
	code := strings.Join(strings.Split(item.CurrentCode, ""), "   ")

	// Render
	fmt.Fprint(w, render_fn(m.Width()-9, account, issuer, code))
}

// Tokens
type Tokens struct {
	// Vault
	vault *tlockvault.Vault

	// Current folder
	folder *tlockvault.Folder

	// Tokens
	listview *list.Model

	// Context
	context *context.Context
}

// Returns the focused folder item
func (tokens Tokens) Focused() *tokensListItem {
	// If there are no items, return nil
	if tokens.listview == nil || len(tokens.listview.Items()) == 0 {
		return nil
	}

	// Get the focused item
	focusedToken := tokens.listview.Items()[tokens.listview.Index()].(tokensListItem)

	// Return
	return &focusedToken
}

// Builds the token list view items
func buildTokensItems(tokens []tlockvault.Token) []list.Item {
	mapper := func(token tlockvault.Token) list.Item {
		return InitializeTokenListItem(token)
	}

	return tlockinternal.Map(tokens, mapper)
}

// Builds the tokens list view
func buildTokensListView(tokens []tlockvault.Token) list.Model {
	// Get terminal size
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	return components.ListViewSimple(buildTokensItems(tokens), tokensListDelegate{}, tokensWidth(width), height-3)
}

// Initializes a new instance of folders
func InitializeTokens(vault *tlockvault.Vault, context *context.Context) Tokens {
	return Tokens{
		vault:   vault,
		folder:  nil,
		context: context,
	}
}

// Handles update messages
func (tokens *Tokens) Update(msg tea.Msg, manager *modelmanager.ModelManager) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msgType.String() == "c":
			if focused := tokens.Focused(); focused != nil && tokens.context.ClipboardAvailability {
				clipboard.Write(clipboard.FmtText, []byte(focused.CurrentCode))
			}

		case msgType.String() == "a":
			if tokens.folder != nil {
				manager.PushScreen(InitializeAddTokenScreen(*tokens.folder, tokens.vault))
			}

		case msgType.String() == "e":
			if focused := tokens.Focused(); focused != nil {
				manager.PushScreen(InitializeEditTokenScreen(*tokens.folder, focused.Token, tokens.vault))
			}

		case msgType.String() == "m":
			if focused := tokens.Focused(); focused != nil {
				manager.PushScreen(InitializeMoveTokenScreen(tokens.vault, *tokens.folder, focused.Token))
			}

		case msgType.String() == "d":
			if focused := tokens.Focused(); focused != nil {
				manager.PushScreen(InitializeDeleteTokenScreen(tokens.vault, *tokens.folder, focused.Token))
			}

		case msgType.String() == "J":
			if focused := tokens.Focused(); focused != nil {
				// Move token down
				tokens.vault.MoveTokenDown(tokens.folder.ID, focused.Token.ID)

				// Refresh tokens
				cmds = append(cmds, func() tea.Msg { return tlockmessages.RefreshTokensMsg{} })

				// Move cursor down
				tokens.listview.CursorDown()
			}

		case msgType.String() == "K":
			if focused := tokens.Focused(); focused != nil {
				// Move token down
				tokens.vault.MoveTokenUp(tokens.folder.ID, focused.Token.ID)

				// Refresh tokens
				cmds = append(cmds, func() tea.Msg { return tlockmessages.RefreshTokensMsg{} })

				// Move cursor down
				tokens.listview.CursorUp()
			}

		case msgType.String() == "n":
			if focused := tokens.Focused(); focused != nil {
				if focused.Token.Type == tlockvault.TokenTypeHOTP {
					tokens.vault.IncreaseCounter(tokens.folder.ID, focused.Token.ID)

					// Refresh tokens
					cmds = append(cmds, func() tea.Msg { return tlockmessages.RefreshTokensMsg{} })
				}
			}

		case key.Matches(msgType, tokenKeys.Screen):
			if tokens.folder != nil {
				manager.PushScreen(InitializeTokenFromScreen(tokens.vault, *tokens.folder))
			}
		}

	case tlockmessages.FolderChanged:
		// Build listview
		listview := buildTokensListView(tokens.vault.GetTokens(msgType.Folder.ID))

		// Update listview
		tokens.listview = &listview
		tokens.folder = &msgType.Folder

	case tlockmessages.RefreshTokensValue:
		if tokens.listview != nil {
			items := make([]list.Item, len(tokens.listview.Items()))

			for index, item := range tokens.listview.Items() {
				tokenItem := item.(tokensListItem)
				tokenItem.Refresh()

				items[index] = tokenItem
			}

			cmds = append(cmds, tokens.listview.SetItems(items))
		}

	case tea.WindowSizeMsg:
		if tokens.listview != nil {
			tokens.listview.SetWidth(tokensWidth(msgType.Width))
			tokens.listview.SetHeight(msgType.Height - 3)
		}

	case tlockmessages.RefreshTokensMsg:
		if tokens.folder != nil {
			cmds = append(cmds, tokens.listview.SetItems(buildTokensItems(tokens.vault.GetTokens(tokens.folder.ID))))
		}
	}

	// Update listview
	if tokens.listview != nil {
		updatedListView, _ := tokens.listview.Update(msg)
		tokens.listview = &updatedListView
	}

	return tea.Batch(cmds...)
}

// View
func (tokens Tokens) View() string {
	if tokens.folder == nil {
		// Yet to recieve message
		return ""
	}

	// Render placeholder for no tokens
	if len(tokens.listview.Items()) == 0 {
		style := lipgloss.NewStyle().
			Height(tokens.listview.Height()).
			Width(tokens.listview.Width()).
			Align(lipgloss.Center, lipgloss.Center)

		ui := lipgloss.JoinVertical(
			lipgloss.Center,
			tlockstyles.Styles.Title.Render(EmptyAsciiArt),
			tlockstyles.Styles.SubText.Render("So empty! How about adding a new token?"), "",
			tlockstyles.Help.View(tokenKeys),
		)

		return style.Render(ui)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		tlockstyles.Styles.AccentBgItem.Render("TOKENS"), "",
		tokens.listview.View(),
	)
}
