package tokens

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/context"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/messages"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-internal/utils"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	tlockstyles "github.com/eklairs/tlock/tlock/styles"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
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
func getCurrentCode(token tlockvault.Token) string {
	var code string

	if token.Type == tlockvault.TokenTypeTOTP {
		code, _ = totp.GenerateCodeCustom(token.Secret, time.Now(), totp.ValidateOpts{
			Period:    uint(token.Period),
			Digits:    otp.Digits(token.Digits),
			Algorithm: token.HashingAlgorithm,
		})
	} else {
		code, _ = hotp.GenerateCodeCustom(token.Secret, uint64(token.UsageCounter+token.InitialCounter), hotp.ValidateOpts{
			Digits:    otp.Digits(token.Digits),
			Algorithm: token.HashingAlgorithm,
		})
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
	item.CurrentCode = getCurrentCode(item.Token)
}

// Initializes a new instance of the tokens list item
func InitializeTokenListItem(token tlockvault.Token) tokensListItem {
	var ttr *int

	if token.Type == tlockvault.TokenTypeTOTP {
		timeToRefresh := getRemainingTime(token)
		ttr = &timeToRefresh
	}

	return tokensListItem{
		CurrentCode: getCurrentCode(token),
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
var tokenKeys tokenKeyMap

var EmptyAsciiArt = `
\    /\
 )  ( ')
(  /  )
 \(__)|
`

// Tokens list delegate
type tokensListDelegate struct {
	context *context.Context
}

// Height
func (d tokensListDelegate) Height() int {
	return 3
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
	render_fn := components.TokenItemActive

	if index != m.Index() {
		render_fn = components.TokenItemInactive
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
	codeToShow := item.CurrentCode

	// If it is not the current focused index, we will hide the token
	if index != m.Index() {
		// Hide it aywa with astrisk
		codeToShow = strings.Repeat("*", item.Token.Digits)
	}

	var tokenRenderable string

	icon, ok := d.context.Icons[item.Token.Issuer]

	if ok {
		tokenRenderable = lipgloss.NewStyle().Foreground(lipgloss.Color("#" + icon.Hex)).Bold(true).Render(icon.Unicode)
	} else {
		tokenRenderable = tlockstyles.Styles.Title.Render("")
	}

	// Render
	fmt.Fprint(w, render_fn(m.Width()-9, tokenRenderable, account, issuer, strings.Join(strings.Split(codeToShow, ""), "   "), item.Token.Period, item.time, d.context.Config.EnableIcons))
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

	return utils.Map(tokens, mapper)
}

// Builds the tokens list view
func buildTokensListView(tokens []tlockvault.Token, context *context.Context) list.Model {
	// Get terminal size
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	return components.ListViewSimple(buildTokensItems(tokens), tokensListDelegate{context: context}, tokensWidth(width), height-5)
}

// Initializes a new instance of folders
func InitializeTokens(vault *tlockvault.Vault, context *context.Context) Tokens {
	// Initialize keys
	tokenKeys = tokenKeyMap{
		Manual: key.NewBinding(
			key.WithKeys(context.Config.Tokens.Add.Keys()...),
			key.WithHelp(strings.Join(context.Config.Tokens.Add.Keys(), "/"), "add token"),
		),
		Screen: key.NewBinding(
			key.WithKeys(context.Config.Tokens.AddScreen.Keys()...),
			key.WithHelp(strings.Join(context.Config.Tokens.AddScreen.Keys(), "/"), "add token from screen"),
		),
	}

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
		case key.Matches(msgType, tokens.context.Config.Tokens.Copy.Binding):
			if clipboard.Unsupported {
				cmds = append(cmds, func() tea.Msg {
					return components.StatusBarMsg{Message: "Clipboard is not available", ErrorMessage: true}
				})
			} else {
				if focused := tokens.Focused(); focused != nil {
					// Set clipboard
					clipboard.WriteAll(focused.CurrentCode)

					accountName := focused.Token.Account

					if accountName == "" {
						accountName = "<no account name>"
					}

					cmds = append(cmds, func() tea.Msg {
						return components.StatusBarMsg{Message: fmt.Sprintf("Successfully copied token (%s)", accountName)}
					})
				}

			}

		case key.Matches(msgType, tokens.context.Config.Tokens.Add.Binding):
			if tokens.folder != nil {
				manager.PushScreen(InitializeAddTokenScreen(*tokens.folder, tokens.vault))
			}

		case key.Matches(msgType, tokens.context.Config.Tokens.Edit.Binding):
			if focused := tokens.Focused(); focused != nil {
				manager.PushScreen(InitializeEditTokenScreen(*tokens.folder, focused.Token, tokens.vault))
			}

		case key.Matches(msgType, tokens.context.Config.Tokens.Move.Binding):
			if focused := tokens.Focused(); focused != nil {
				manager.PushScreen(InitializeMoveTokenScreen(tokens.vault, *tokens.folder, focused.Token))
			}

		case key.Matches(msgType, tokens.context.Config.Tokens.Delete.Binding):
			if focused := tokens.Focused(); focused != nil {
				manager.PushScreen(InitializeDeleteTokenScreen(tokens.vault, *tokens.folder, focused.Token))
			}

		case key.Matches(msgType, tokens.context.Config.Tokens.MoveDown.Binding):
			if focused := tokens.Focused(); focused != nil {
				// Move token down
				tokens.vault.MoveTokenDown(tokens.folder.Name, focused.Token)

				// Refresh tokens
				cmds = append(cmds, func() tea.Msg { return tlockmessages.RefreshTokensMsg{} })

				// Move cursor down
				tokens.listview.CursorDown()

				accountName := focused.Token.Account

				if accountName == "" {
					accountName = "<no account name>"
				}

				cmds = append(cmds, func() tea.Msg {
					return components.StatusBarMsg{Message: fmt.Sprintf("Successfully moved the token down (%s)", accountName)}
				})
			}

		case key.Matches(msgType, tokens.context.Config.Tokens.MoveUp.Binding):
			if focused := tokens.Focused(); focused != nil {
				// Move token down
				tokens.vault.MoveTokenUp(tokens.folder.Name, focused.Token)

				// Refresh tokens
				cmds = append(cmds, func() tea.Msg { return tlockmessages.RefreshTokensMsg{} })

				// Move cursor down
				tokens.listview.CursorUp()

				accountName := focused.Token.Account

				if accountName == "" {
					accountName = "<no account name>"
				}

				cmds = append(cmds, func() tea.Msg {
					return components.StatusBarMsg{Message: fmt.Sprintf("Successfully moved the token up (%s)", accountName)}
				})
			}

		case key.Matches(msgType, tokens.context.Config.Tokens.NextHOTP.Binding):
			if focused := tokens.Focused(); focused != nil {
				if focused.Token.Type == tlockvault.TokenTypeHOTP {
					tokens.vault.IncreaseCounter(tokens.folder.Name, focused.Token)

					accountName := focused.Token.Account
					if accountName == "" {
						accountName = "<no account name>"
					}

					// Refresh tokens
					cmds = append(cmds, func() tea.Msg { return tlockmessages.RefreshTokensMsg{} })
					cmds = append(cmds, func() tea.Msg {
						return components.StatusBarMsg{Message: fmt.Sprintf("Successfully generated next token for %s", accountName)}
					})
				}
			}

		case key.Matches(msgType, tokens.context.Config.Tokens.AddScreen.Binding):
			if tokens.folder != nil {
				cmds = append(cmds, manager.PushScreen(InitializeTokenFromScreen(tokens.vault, *tokens.folder)))
			}
		}

	case tlockmessages.FolderChanged:
		// Build listview
		listview := buildTokensListView(tokens.vault.GetTokens(msgType.Folder.Name), tokens.context)

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
			tokens.listview.SetHeight(msgType.Height - 5)
		}

	case tlockmessages.RefreshTokensMsg:
		if tokens.folder != nil {
			cmds = append(cmds, tokens.listview.SetItems(buildTokensItems(tokens.vault.GetTokens(tokens.folder.Name))))
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
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))

	if tokens.folder == nil {
		// Yet to recieve message
		return ""
	}

	// Render placeholder for no tokens
	if len(tokens.listview.Items()) == 0 {
		style := lipgloss.NewStyle().
			Height(height-2).
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
		lipgloss.Left, "",
		tlockstyles.Styles.AccentBgItem.Render("TOKENS"), "",
		tokens.listview.View(),
	)
}
