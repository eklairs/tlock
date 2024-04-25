package tokens

import (
	"math"
	"os"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/term"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-models/dashboard/folders"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

// Returns the current code
func getCurrentCode(key *otp.Key, usageCounter int) string {
	var code string

	if key.Type() == "totp" {
		code, _ = totp.GenerateCode(key.Secret(), time.Now())
	} else {
		code, _ = hotp.GenerateCode(key.Secret(), uint64(usageCounter))
	}

	return code
}

// Token list item
type tokensListItem struct {
	// Type of token
	TokenType string

	// Current code
	CurrentCode string

	// URI string
	Token tlockvault.Token
}

// Initializes a new instance of the tokens list item
func InitializeTokenListItem(token tlockvault.Token) tokensListItem {
	tokenKey, _ := otp.NewKeyFromURL(token.URI)

	return tokensListItem{
		TokenType:   tokenKey.Type(),
		CurrentCode: getCurrentCode(tokenKey, token.UsageCounter),
		Token:       token,
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

// Tokens
type Tokens struct {
	// Vault
	vault *tlockvault.TLockVault

	// Folder
	folder *string

	// Help
	help help.Model

	// Tokens
	tokens []tlockvault.Token
}

// Initializes a new instance of folders
func InitializeTokens(vault *tlockvault.TLockVault) Tokens {
	return Tokens{
		vault:  vault,
		folder: nil,
		help:   components.BuildHelp(),
	}
}

// Init
func (tokens *Tokens) Init() tea.Cmd {
	return nil
}

// Handles update messages
func (tokens *Tokens) Update(msg tea.Msg, manager *modelmanager.ModelManager) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case folders.FolderChangedMsg:
		tokens.folder = &msgType.Folder
		tokens.tokens = tokens.vault.GetTokens(msgType.Folder)

	case tea.KeyMsg:
		switch msgType.String() {
		case "a":
			manager.PushScreen(InitializeAddTokenScreen())
		case "s":
			manager.PushScreen(InitializeTokenFromScreen())
		}
	}

	return tea.Batch(cmds...)
}

// View
func (tokens Tokens) View() string {
	// Get size
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	// Width
	tokensWidth := width - int(math.Floor((1.0/5.0)*float64(width)))

	// If no tokens, render placeholder
	if len(tokens.tokens) == 0 {
		style := lipgloss.NewStyle().
			Height(height-3).
			Width(tokensWidth).
			Align(lipgloss.Center, lipgloss.Center)

		ui := lipgloss.JoinVertical(
			lipgloss.Center,
			tlockstyles.Styles.Title.Render(EmptyAsciiArt),
			tlockstyles.Styles.SubText.Render("So empty! How about adding a new token?"), "",
			tokens.help.View(tokenKeys),
		)

		return style.Render(ui)
	}

	return ""
}
