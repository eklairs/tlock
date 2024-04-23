package tokens

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/term"

	"github.com/eklairs/tlock/tlock-internal/boundedinteger"
	"github.com/eklairs/tlock/tlock-internal/buildhelp"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-models/dashboard/folders"

	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

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
	// Context
	context context.Context

	// Vault
	vault *tlockvault.TLockVault

	// Focused index
	focused_index boundedinteger.BoundedInteger

	// Styles
	styles tlockstyles.Styles

	// Folder
	folder *string

	// Help
	help help.Model

	// Tokens
	tokens []tlockvault.TokenURI
}

// Initializes a new instance of folders
func InitializeTokens(vault *tlockvault.TLockVault, context context.Context) Tokens {
	// Terminal size
	width, _, _ := term.GetSize(0)

	// Styles
	styles := tlockstyles.InitializeStyle(width-folders.FOLDERS_WIDTH-6, context.Theme)

	return Tokens{
		vault:   vault,
		styles:  styles,
		context: context,
		folder:  nil,
		help:    buildhelp.BuildHelp(styles),
	}
}

type refreshTokens struct{}

func notifyRefreshTokens() tea.Msg {
    time.Sleep(time.Second * 1)

    return refreshTokens{}
}

func (tokens *Tokens) Init() tea.Cmd {
    return notifyRefreshTokens
}

// Handles update messages
func (tokens *Tokens) Update(msg tea.Msg, manager *modelmanager.ModelManager) tea.Cmd {
    cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "j":
			tokens.focused_index.Increase()
		case "k":
			tokens.focused_index.Decrease()
		case "s":
			manager.PushScreen(InitializeTokenFromScreen(tokens.context))
        case "a":
            manager.PushScreen(InitializeAddTokenModel())
        case "e":
            focused_uri := tokens.vault.GetTokens(*tokens.folder)[tokens.focused_index.Value].URI

            manager.PushScreen(InitializeEditTokenModel(focused_uri))
        case "n":
            focused_uri := tokens.vault.GetTokens(*tokens.folder)[tokens.focused_index.Value].URI

            authKey, _ := otp.NewKeyFromURL(focused_uri)

            if authKey.Type() == "hotp" {
                tokens.vault.IncrementTokenUsageCounter(*tokens.folder, focused_uri)
            }
		}

	case AddTokenMsg:
		if tokens.folder != nil {
			tokens.vault.AddToken(*tokens.folder, msgType.URI)

            tokens.tokens = tokens.vault.GetTokens(*tokens.folder)
            tokens.focused_index = boundedinteger.New(tokens.focused_index.Value, len(tokens.tokens))
		}

	case EditTokenMsg:
		if tokens.folder != nil {
            tokens.vault.EditToken(*tokens.folder, msgType.Old, msgType.New)

            tokens.tokens = tokens.vault.GetTokens(*tokens.folder)
            tokens.focused_index = boundedinteger.New(tokens.focused_index.Value, len(tokens.tokens))
		}

	case folders.FolderChangedMsg:
		tokens.folder = &msgType.Folder
		tokens.tokens = tokens.vault.GetTokens(msgType.Folder)

		tokens.focused_index = boundedinteger.New(0, len(tokens.tokens))

    case refreshTokens:
        cmds = append(cmds, notifyRefreshTokens)
	}

	return tea.Batch(cmds...)
}

// View
func (tokens Tokens) View() string {
	// Get term size
	width, height, _ := term.GetSize(0)

	// If the folder is nil, then we have not yet recieved the change folder message, so we should not render anything
	if tokens.folder == nil {
		return ""
	}

    tokensWidth := width - folders.FOLDERS_WIDTH - 13

	if len(tokens.tokens) == 0 {
		style := tokens.styles.Base.Copy().
			Height(height - 3).
			Align(lipgloss.Center, lipgloss.Center)

		ui := lipgloss.JoinVertical(
			lipgloss.Left,
			tokens.styles.Center.Render(tokens.styles.Title.Copy().UnsetWidth().Render(EmptyAsciiArt)),
			tokens.styles.Center.Render(tokens.styles.Title.Copy().UnsetBold().UnsetWidth().Render("So empty! How about adding a new token?")), "",
			tokens.styles.Center.Copy().UnsetWidth().Render(tokens.help.View(tokenKeys)),
		)

		return style.Render(ui)
	}

	// List of items
	items := make([]string, 0)

    // Header
    items = append(items, tokens.styles.AccentTitle.Copy().Margin(1).Render("TOKENS"))

	// Iter
	for index, token := range tokens.tokens {
		// Generate key
		authKey, _ := otp.NewKeyFromURL(token.URI)

        // Token info
        tokenInfo := lipgloss.JoinHorizontal(
            lipgloss.Left,
            tokens.styles.Dimmed.Copy().UnsetWidth().Render(authKey.AccountName()),
            tokens.styles.Dimmed.Copy().UnsetWidth().Render(" • "),
            tokens.styles.Dimmed.Copy().UnsetWidth().Render(authKey.Issuer()),
        )

        // Render fn
        render_fn := tokens.styles.InactiveListItem.Render

		// Check if it is the focused one
		if index == tokens.focused_index.Value {
			tokenInfo = lipgloss.JoinHorizontal(
				lipgloss.Left,
				tokens.styles.Title.Copy().UnsetWidth().Render(authKey.AccountName()),
				tokens.styles.Dimmed.Copy().UnsetWidth().Background(tokens.context.Theme.WindowBgOver).Render(" • "),
				tokens.styles.Base.Copy().UnsetWidth().Background(tokens.context.Theme.WindowBgOver).Render(authKey.Issuer()),
			)

            render_fn = tokens.styles.ActiveItem.Render
		}

        // Token str
        var currentToken string

        if authKey.Type() == "hotp" {
            currentToken, _ = hotp.GenerateCode(authKey.Secret(), uint64(token.UsageCounter))
        } else {
            currentToken, _ = totp.GenerateCode(authKey.Secret(), time.Now())
        }

        tokenRenderable := ""

        style := lipgloss.NewStyle().
            Bold(true).
            PaddingRight(3).
            Background(tokens.context.Theme.WindowBgOver).
            Foreground(tokens.context.Theme.Accent)

        spacer_style := lipgloss.NewStyle().
            Background(tokens.context.Theme.WindowBgOver)

        for _, tokenChar := range currentToken {
            tokenRenderable = lipgloss.JoinHorizontal(lipgloss.Left, tokenRenderable, style.Render(string(tokenChar)))
        }

        timeRemainingRenderable := ""

        // Reset information
        if authKey.Type() == "totp" {
            // https://github.com/pyauth/pyotp/issues/87#issuecomment-561284149
            timeRemaining := authKey.Period() - uint64(time.Now().Unix()) % authKey.Period()

            timeRemainingRenderable = tokens.styles.Dimmed.Copy().UnsetWidth().Italic(true).Render(fmt.Sprintf("Refreshing in %d", timeRemaining))
        }

        // Final ui
        ui := lipgloss.JoinHorizontal(
            lipgloss.Left,
            tokenInfo,
            spacer_style.Render(strings.Repeat(" ", tokensWidth - lipgloss.Width(tokenInfo) - lipgloss.Width(tokenRenderable) - lipgloss.Width(timeRemainingRenderable))),
            tokenRenderable,
            timeRemainingRenderable,
        )

		// items = append(items, render_fn(ui))
        items = append(items, render_fn(ui))
	}

	// Render
	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

