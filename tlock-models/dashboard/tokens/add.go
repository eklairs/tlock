package tokens

import (
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tlockinternal "github.com/eklairs/tlock/tlock-internal"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/form"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/tlock-messages"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"github.com/google/uuid"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/term"
)

// Error for empty secret
var ERROR_EMPTY_SECRET = "Secret is required"

// Error for invalid secret
var ERROR_INVALID_SECRET = "Secret is invalid, are you sure you have typed correctly?"

// Add token key map
type addTokenKeyMap struct {
	GoBack key.Binding
	Enter  key.Binding
	Tab    key.Binding
	Arrow  key.Binding
}

// ShortHelp()
func (k addTokenKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Arrow, k.Enter, k.GoBack}
}

// FullHelp()
func (k addTokenKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

// Keys
var addTokenKeys = addTokenKeyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "create token"),
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

var addTokenAscii = `
▄▀█ █▀▄ █▀▄
█▀█ █▄▀ █▄▀`

// Add token screen
type AddTokenScreen struct {
	// Form
	form form.Form

	// Vault
	vault *tlockvault.Vault

	// Folder
	folder tlockvault.Folder

	// Viewport
	viewport viewport.Model

	// Viewport content
	content string
}

// Initializes a new screen of AddTokenScreen
func InitializeAddTokenScreen(folder tlockvault.Folder, vault *tlockvault.Vault) AddTokenScreen {
	items := []form.FormItem{
		form.FormItemInputBox{
			Title:       "Account Name",
			Description: "Name of the account, like John Doe",
			Input:       components.InitializeInputBox("Account name goes here..."),
		},
		form.FormItemInputBox{
			Title:       "Issuer Name",
			Description: "Name of the issuer, like GitHub",
			Input:       components.InitializeInputBox("Issuer name goes here..."),
		},
		form.FormItemInputBox{
			Title:       "Secret",
			Description: "The secret provided by the issuer",
			Input:       components.InitializeInputBox("The secret goes here..."),
		},
		form.FormItemOptionBox{
			Title:         "Type",
			Description:   "Time or counter based token",
			Values:        []string{"TOTP", "HOTP"},
			SelectedIndex: 0,
		},
		form.FormItemOptionBox{
			Title:         "Hash function",
			Description:   "Hash function for the token",
			Values:        []string{"SHA1", "SHA256", "SHA512", "MD5"},
			SelectedIndex: 1,
		},
		form.FormItemInputBox{
			Title:       "Period",
			Description: "Time to refresh the token",
			Input:       tlockinternal.ValidatorInteger(components.InitializeInputBoxCustomWidth("Time in seconds...", 24)),
		},
		form.FormItemInputBox{
			Title:       "Initial counter",
			Description: "Initial counter for HOTP token",
			Input:       tlockinternal.ValidatorInteger(components.InitializeInputBoxCustomWidth("Initial counter...", 24)),
		},
		form.FormItemInputBox{
			Title:       "Digits",
			Description: "Number of digits",
			Input:       tlockinternal.ValidatorInteger(components.InitializeInputBoxCustomWidth("Number of digits goes here...", 24)),
		},
	}

	// Disable the HOTP initially
	form := form.New(items)
	form.Items[6].Enabled = false

	// Get term size
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))

	// Initialize viewport
	content := GenerateUI(form)

	viewport := viewport.New(85, min(height, lipgloss.Height(content)))
	viewport.SetContent(content)

	// Restrict to arrow keys to move viewport up and down to prevent conflict when the user is typing
	viewport.KeyMap.Up = key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "up"),
	)

	viewport.KeyMap.Down = key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "down"),
	)

	// Return
	return AddTokenScreen{
		form:     form,
		vault:    vault,
		folder:   folder,
		content:  content,
		viewport: viewport,
	}
}

// Init
func (screen AddTokenScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen AddTokenScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	// Get the secret item
	secretItem := screen.form.Items[2].FormItem.(form.FormItemInputBox)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, addTokenKeys.Enter):
			// Get the user secret
			secret := secretItem.Value()

			// Error if any
			var error *string

			if secret == "" {
				error = &ERROR_EMPTY_SECRET
			}

			// Try to generate code with the secret
			_, err := totp.GenerateCode(secret, time.Now())

			if err != nil {
				error = &ERROR_INVALID_SECRET
			}

			if error != nil {
				// Set the error
				secretItem.ErrorMessage = error

				// Update item
				screen.form.Items[2].FormItem = secretItem

				// Break
				break
			}

			// Instance of form
			formItems := screen.form.Items

			// Type of token
			var tokenType tlockvault.TokenType

			if formItems[3].FormItem.Value() == "TOTP" {
				tokenType = tlockvault.TokenTypeTOTP
			} else {
				tokenType = tlockvault.TokenTypeHOTP
			}

			// Hashing function
			hashFunction := otp.AlgorithmSHA256

			switch formItems[4].FormItem.Value() {
			case "SHA1":
				hashFunction = otp.AlgorithmSHA1
			case "SHA512":
				hashFunction = otp.AlgorithmSHA512
			case "MD5":
				hashFunction = otp.AlgorithmMD5
			}

			// To int
			toInt := func(str string) int {
				res, _ := strconv.Atoi(str)

				return res
			}

			// This or that
			or := func(left string, right int) int {
				if left == "" {
					return right
				}

				return toInt(left)
			}

			// Okay its time to add!
			token := tlockvault.Token{
				ID:               uuid.NewString(),
				Account:          formItems[0].FormItem.Value(),
				Issuer:           formItems[1].FormItem.Value(),
				Secret:           formItems[2].FormItem.Value(),
				Type:             tokenType,
				HashingAlgorithm: hashFunction,
				Period:           or(formItems[5].FormItem.Value(), 30),
				InitialCounter:   or(formItems[6].FormItem.Value(), 0),
				Digits:           or(formItems[7].FormItem.Value(), 6),
				UsageCounter:     0,
			}

			// Add
			screen.vault.AddTokenFromToken(screen.folder.ID, token)

			// Require refresh of folders and tokens list
			cmds = append(
				cmds,
				func() tea.Msg { return tlockmessages.RefreshFoldersMsg{} },
				func() tea.Msg { return tlockmessages.RefreshTokensMsg{} },
			)

			// Pop
			manager.PopScreen()

		case key.Matches(msgType, editTokenKeys.GoBack):
			manager.PopScreen()
		}

	case tea.WindowSizeMsg:
		screen.viewport.Height = min(msgType.Height, lipgloss.Height(screen.content))
	}

	// Update viewport
	screen.viewport, _ = screen.viewport.Update(msg)

	// Update the form
	screen.form.Update(msg)

	// Generate UI
	screen.content = GenerateUI(screen.form)

	// Set viewport content
	screen.viewport.SetContent(screen.content)

	// Enable / Disable items based on the choosen type
	if screen.form.Items[3].FormItem.Value() == "TOTP" {
		// Enable period item
		screen.form.Items[5].Enabled = true

		// Disable initial counter item
		screen.form.Items[6].Enabled = false
	}

	if screen.form.Items[3].FormItem.Value() == "HOTP" {
		// Enable initial counter item
		screen.form.Items[6].Enabled = true

		// Disable period item
		screen.form.Items[5].Enabled = false
	}

	return screen, tea.Batch(cmds...)
}

// View
func (screen AddTokenScreen) View() string {
	return screen.viewport.View()
}

// Generates the UI
func GenerateUI(form form.Form) string {
	// Items
	items := []string{
		tlockstyles.Styles.Title.Render(addTokenAscii), "",
		tlockstyles.Styles.SubText.Render("Add a new token [secret is required, rest are optional]"), "",
		form.Items[0].FormItem.View(), // Account name input
		form.Items[1].FormItem.View(), // Issuer name input
		form.Items[2].FormItem.View(), // Secret value input
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			form.Items[3].FormItem.View(), "   ",
			form.Items[4].FormItem.View(),
		), "",
	}

	// Render the input boxes based on choosen type
	inputGroup := lipgloss.JoinHorizontal(
		lipgloss.Left,
		form.Items[5].FormItem.View(), "   ",
		form.Items[7].FormItem.View(),
	)

	if form.Items[3].FormItem.Value() == "HOTP" {
		inputGroup = lipgloss.JoinHorizontal(
			lipgloss.Left,
			form.Items[6].FormItem.View(), "   ",
			form.Items[7].FormItem.View(),
		)
	}

	// Add the help menu
	items = append(items, inputGroup, "", tlockstyles.Help.View(addTokenKeys))

	// Return!
	return lipgloss.JoinVertical(lipgloss.Center, items...)
}
