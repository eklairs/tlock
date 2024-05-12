package tokens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/form"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/messages"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"github.com/pquerna/otp"
)

func value(input textinput.Model, value string) textinput.Model {
	input.SetValue(value)

	return input
}

// Edit token key map
type editTokenKeyMap struct {
	GoBack key.Binding
	Enter  key.Binding
	Tab    key.Binding
	Arrow  key.Binding
}

// ShortHelp()
func (k editTokenKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Arrow, k.Enter, k.GoBack}
}

// FullHelp()
func (k editTokenKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

// Keys
var editTokenKeys = editTokenKeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "edit token"),
	),
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab", "shift+tab"),
		key.WithHelp("tab/shift+tab", "switch input"),
	),
	Arrow: key.NewBinding(
		key.WithKeys("right", "left"),
		key.WithHelp("→/←", "change option"),
	),
}

// Edit token ascii art
var editTokenAscii = `
█▀▀ █▀▄ █ ▀█▀
██▄ █▄▀ █  █`

// Edit token desc
var editTokenDesc = "Edit your token [secret is required, rest are optional]"

// Edit token screen
type EditTokenScreen struct {
	// Form
	form form.Form

	// Vault
	vault *tlockvault.Vault

	// Folder
	folder tlockvault.Folder

	// Token to edit
	token tlockvault.Token

	// Viewport
	viewport viewport.Model

	// Viewport content
	content string
}

func tokenTypeToString(tokentype tlockvault.TokenType) string {
    if tokentype == tlockvault.TokenTypeHOTP {
        return "HOTP"
    }

    return "TOTP"
}

func hashAlgoToString(hash otp.Algorithm) string {
    switch hash {
    case otp.AlgorithmSHA256:
        return "SHA256"
    case otp.AlgorithmSHA512:
        return "SHA512"
    }

    return "SHA1"
}

// Initializes a new screen of EditTokenScreen
func InitializeEditTokenScreen(folder tlockvault.Folder, token tlockvault.Token, vault *tlockvault.Vault) EditTokenScreen {
    // Form
    form := BuildForm(map[string]string {
        "account": token.Account,
        "issuer": token.Issuer,
        "secret": token.Secret,
        "type": tokenTypeToString(token.Type),
        "hash": hashAlgoToString(token.HashingAlgorithm),
        "period": fmt.Sprintf("%d", token.Period),
        "digits": fmt.Sprintf("%d", token.Digits),
        "counter": fmt.Sprintf("%d", token.InitialCounter),
    })

    // Override the validator for edit screen for secret
    form.Items[2].Validators = []func(vault *tlockvault.Vault, value string) error{
        func(vault *tlockvault.Vault, secret string) error {
            // Validate
            _, err := vault.ValidateToken(secret)

            // If it is error because of duplicate secret, and if the secret is same, lets ignore it
            if secret == token.Secret && err == tlockvault.ERR_TOKEN_EXISTS {
                return nil
            }

            // Return
            return err
        },
    }

    // Return
	return EditTokenScreen{
		form:     form,
		token:    token,
		vault:    vault,
		folder:   folder,
		viewport: IntoViewport(editTokenAscii, editTokenDesc, form),
	}
}

// Init
func (screen EditTokenScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen EditTokenScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
    // Commands
    cmds := make([]tea.Cmd, 0)

    // Match
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msgType, addTokenKeys.GoBack):
            manager.PopScreen()
        }

    case form.FormSubmittedMsg:
        // Get token
        token := TokenFromFormData(msgType.Data)

        // Make statusbar message
        statusBarMessage := fmt.Sprintf("Successfully edited token for %s", screen.token.Account)

        if screen.token.Account == "" {
            statusBarMessage = fmt.Sprintf("Successfully edited token (no account name)")
        }

        // Require refresh of folders and tokens list
        cmds = append(
            cmds,
            func() tea.Msg { return tlockmessages.RefreshFoldersMsg{} },
            func() tea.Msg { return tlockmessages.RefreshTokensMsg{} },
            func() tea.Msg { return components.StatusBarMsg{Message: statusBarMessage} },
        )

        // Add
        screen.vault.ReplaceToken(screen.folder.Name, screen.token, token)

        // Break
        manager.PopScreen()
    }

    // Let the form handle
    cmds = append(cmds, screen.form.Update(msg, screen.vault))

    // Update the viewport
    DisableBasedOnType(&screen.form)
    UpdateViewport(editTokenAscii, editTokenDesc, &screen.viewport, screen.form)

    // Return
	return screen, tea.Batch(cmds...)
}

// View
func (screen EditTokenScreen) View() string {
	return screen.viewport.View()
}

