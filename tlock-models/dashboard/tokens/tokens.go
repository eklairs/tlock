package tokens

import (
	"io"
	"math"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tlockinternal "github.com/eklairs/tlock/tlock-internal"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/tlock-messages"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/term"
)

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

}

// Tokens
type Tokens struct {
	// Vault
	vault *tlockvault.Vault

	// Current folder
	folder *tlockvault.Folder

	// Tokens
	tokensListView *list.Model
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

	// Calculate width
	tokensWidth := width - int(math.Floor((1.0/5.0)*float64(width)))

	return components.ListViewSimple(buildTokensItems(tokens), tokensListDelegate{}, tokensWidth, height)
}

// Initializes a new instance of folders
func InitializeTokens(vault *tlockvault.Vault) Tokens {
	return Tokens{
		vault:  vault,
		folder: nil,
	}
}

// Handles update messages
func (tokens *Tokens) Update(msg tea.Msg, manager *modelmanager.ModelManager) tea.Cmd {
	switch msgType := msg.(type) {
	case tlockmessages.FolderChanged:
		// Build listview
		listview := buildTokensListView(tokens.vault.GetTokens(msgType.Folder.ID))

		// Update listview
		tokens.tokensListView = &listview
		tokens.folder = &msgType.Folder
	}

	return nil
}

// View
func (tokens Tokens) View() string {
	if tokens.folder == nil {
		// Yet to recieve message
		return ""
	}

	// Render placeholder for no tokens
	if len(tokens.tokensListView.Items()) == 0 {
		style := lipgloss.NewStyle().
			Height(tokens.tokensListView.Height()).
			Width(tokens.tokensListView.Width()).
			Align(lipgloss.Center, lipgloss.Center)

		ui := lipgloss.JoinVertical(
			lipgloss.Center,
			tlockstyles.Styles.Title.Render(EmptyAsciiArt),
			tlockstyles.Styles.SubText.Render("So empty! How about adding a new token?"), "",
			tlockstyles.Help.View(tokenKeys),
		)

		return style.Render(ui)
	}

	return ""
}
