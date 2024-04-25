package tokens

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
	"golang.design/x/clipboard"
	"golang.org/x/term"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/context"
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

type secondPassedMsg struct{}

// Token list item
type tokensListItem struct {
	// Type of token
	TokenType string

	// Current code
	CurrentCode string

	// URI string
	Token tlockvault.Token

	// Time remaining before the otp is updated
	// Only in case of totp tokens
	time *int
}

// Refreshes the token
func (item *tokensListItem) Refresh() {
	tokenKey, _ := otp.NewKeyFromURL(item.Token.URI)

	if tokenKey.Type() == "totp" {
		timeToRefresh := int(tokenKey.Period() - uint64(time.Now().Unix())%tokenKey.Period())

		item.time = &timeToRefresh
	}

	item.CurrentCode = getCurrentCode(tokenKey, item.Token.UsageCounter)
}

// Initializes a new instance of the tokens list item
func InitializeTokenListItem(token tlockvault.Token) tokensListItem {
	tokenKey, _ := otp.NewKeyFromURL(token.URI)

	var ttr *int

	if tokenKey.Type() == "totp" {
		timeToRefresh := int(tokenKey.Period() - uint64(time.Now().Unix())%tokenKey.Period())

		ttr = &timeToRefresh
	}

	return tokensListItem{
		TokenType:   tokenKey.Type(),
		CurrentCode: getCurrentCode(tokenKey, token.UsageCounter),
		Token:       token,
		time:        ttr,
	}
}

func (item tokensListItem) FilterValue() string {
	return ""
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
	render_fn := components.ListItemInactive

	if index == m.Index() {
		render_fn = components.ListItemActive
	}

	key, _ := otp.NewKeyFromURL(item.Token.URI)

	// Build key info
	info := lipgloss.JoinHorizontal(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(key.AccountName()),
		tlockstyles.Styles.SubAltBgItem.Render("•"),
		tlockstyles.Styles.SubAltBg.Render(key.Issuer()),
	)

	// Render it differently if it is not the current token
	if index != m.Index() {
		info = lipgloss.JoinHorizontal(
			lipgloss.Center,
			tlockstyles.Styles.SubText.Render(key.AccountName()),
			tlockstyles.Styles.SubText.Render(" • "),
			tlockstyles.Styles.SubText.Render(key.Issuer()),
		)
	}

	// Current code renderable
	current_code := strings.Join(strings.Split(item.CurrentCode, ""), "   ")

	if key.Type() == "totp" {
		timeLeftRenderable := tlockstyles.Styles.SubAltBg.Render(fmt.Sprintf("   ⏲  %d", *item.time))

		if index != m.Index() {
			timeLeftRenderable = tlockstyles.Styles.SubText.Render(fmt.Sprintf("   ⏲  %d", *item.time))
		}

		suffix := lipgloss.JoinHorizontal(
			lipgloss.Center,
			current_code,
			timeLeftRenderable,
		)

		fmt.Fprint(w, render_fn(m.Width()-9, info, suffix))
	} else {
		fmt.Fprint(w, render_fn(m.Width()-9, info, current_code))
	}
}

// Tokens
type Tokens struct {
	// Vault
	vault *tlockvault.TLockVault

	// Folder
	folder *string

	// Tokens
	tokensListView *list.Model

	// Context
	context context.Context
}

// Builds the tokens list view
func buildTokensListView(tokens []tlockvault.Token) list.Model {
	items := make([]list.Item, len(tokens))

	for index, token := range tokens {
		items[index] = InitializeTokenListItem(token)
	}

	return components.ListViewSimple(items, tokensListDelegate{}, 20, 10)
}

// Initializes a new instance of folders
func InitializeTokens(vault *tlockvault.TLockVault, context context.Context) Tokens {
	return Tokens{
		vault:   vault,
		folder:  nil,
		context: context,
	}
}

func notifySecondPassed() tea.Msg {
	currentTime := time.Now()
	nextSecond := currentTime.Truncate(time.Second).Add(time.Second)
	duration := nextSecond.Sub(currentTime)

	time.Sleep(duration)

	return secondPassedMsg{}
}

// Init
func (tokens *Tokens) Init() tea.Cmd {
	return notifySecondPassed
}

// Handles update messages
func (tokens *Tokens) Update(msg tea.Msg, manager *modelmanager.ModelManager) tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case folders.FolderChangedMsg:
		tokens.folder = &msgType.Folder

		tokenListView := buildTokensListView(tokens.vault.GetTokens(msgType.Folder))

		tokens.tokensListView = &tokenListView

	case tea.KeyMsg:
		switch msgType.String() {
		case "a":
			cmds = append(cmds, manager.PushScreen(InitializeAddTokenScreen()))
		case "s":
			cmds = append(cmds, manager.PushScreen(InitializeTokenFromScreen()))
		case "c":
			if tokens.context.ClipboardAvailability {
				item := tokens.tokensListView.Items()[tokens.tokensListView.Index()].(tokensListItem)

				clipboard.Write(clipboard.FmtText, []byte(item.CurrentCode))
			}
        case "m":
            if tokens.folder != nil {
                item := tokens.tokensListView.Items()[tokens.tokensListView.Index()].(tokensListItem)

                cmds = append(cmds, manager.PushScreen(InitializeMoveTokenScreen(tokens.vault, *tokens.folder, item.Token.URI)))
            }
        case "d":
            if tokens.folder != nil {
                item := tokens.tokensListView.Items()[tokens.tokensListView.Index()].(tokensListItem)

                cmds = append(cmds, manager.PushScreen(InitializeDeleteTokenScreen(*tokens.folder, item.Token.URI)))
            }
		}
	case secondPassedMsg:
		if tokens.tokensListView != nil {
			items := make([]list.Item, len(tokens.tokensListView.Items()))

			for index, item := range tokens.tokensListView.Items() {
				tokenItem := item.(tokensListItem)
				tokenItem.Refresh()

				items[index] = tokenItem
			}

			cmds = append(cmds, tokens.tokensListView.SetItems(items))
		}

		cmds = append(cmds, notifySecondPassed)

	case AddTokenMessage:
		if tokens.folder != nil {
			tokens.vault.AddToken(*tokens.folder, msgType.Token)
		}
	}

	if tokens.tokensListView != nil {
		listview, _ := tokens.tokensListView.Update(msg)
		tokens.tokensListView = &listview
	}

	return tea.Batch(cmds...)
}

// View
func (tokens Tokens) View() string {
	// Get size
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	// Width
	tokensWidth := width - int(math.Floor((1.0/5.0)*float64(width))) - 2

	// If no tokens, render placeholder
	if tokens.tokensListView == nil || len(tokens.tokensListView.Items()) == 0 {
		style := lipgloss.NewStyle().
			Height(height-3).
			Width(tokensWidth).
			Align(lipgloss.Center, lipgloss.Center)

		ui := lipgloss.JoinVertical(
			lipgloss.Center,
			tlockstyles.Styles.Title.Render(EmptyAsciiArt),
			tlockstyles.Styles.SubText.Render("So empty! How about adding a new token?"), "",
			tlockstyles.Help.View(tokenKeys),
		)

		return style.Render(ui)
	}

	tokens.tokensListView.SetWidth(tokensWidth)

	// Build UI
	ui := lipgloss.JoinVertical(
		lipgloss.Left,
		tlockstyles.Styles.AccentBgItem.Render("TOKENS"), "",
		tokens.tokensListView.View(),
	)

	// Style
	style := lipgloss.NewStyle().Height(height).Width(tokensWidth)

	// Render
	return style.Render(ui)
}
