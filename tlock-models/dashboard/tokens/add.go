package tokens

import (
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/form"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/messages"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-internal/utils"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/term"
)

// Converts the given string type tok token type
func toTokenType(value string) tlockvault.TokenType {
	if value == "HOTP" {
		return tlockvault.TokenTypeHOTP
	}

	return tlockvault.TokenTypeTOTP
}

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
			Input:       utils.ValidatorIntegerNo0(components.InitializeInputBoxCustomWidth("Time in seconds...", 24)),
		},
		form.FormItemInputBox{
			Title:       "Initial counter",
			Description: "Initial counter for HOTP token",
			Input:       utils.ValidatorInteger(components.InitializeInputBoxCustomWidth("Initial counter...", 24)),
		},
		form.FormItemInputBox{
			Title:       "Digits",
			Description: "Number of digits",
			Input:       utils.ValidatorIntegerNo0(components.InitializeInputBoxCustomWidth("Number of digits goes here...", 24)),
		},
	}

	// Disable the HOTP initially
	form := form.New(items)
	form.Items[6].Enabled = false

	// Get term size
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))

	// Initialize viewport
	content := GenerateUI(form)

	// Initialize viewport and disable conflicting keys
	viewport := utils.DisableViewportKeys(viewport.New(85, min(height, lipgloss.Height(content))))
	viewport.SetContent(content)

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
			for index, item := range screen.form.Items {
				if input, ok := item.FormItem.(form.FormItemInputBox); ok {
					input.ErrorMessage = nil
					screen.form.Items[index].FormItem = input
				}
			}

			// Get the user secret
			secret := secretItem.Value()

			// Dont allow empty secrets
			if secret == "" {
				screen.SetError(2, ERROR_EMPTY_SECRET)
				break
			}

			// Try to generate code with the secret
			_, err := totp.GenerateCode(secret, time.Now())

			if err != nil {
				screen.SetError(2, ERROR_INVALID_SECRET)
				break
			}

			period := utils.Or(screen.form.Items[5].FormItem.Value(), 30)
			digits := utils.Or(screen.form.Items[7].FormItem.Value(), 6)

			// Dont allow period to be zero
			if period < 1 {
				screen.SetError(5, "Period cannot be less than 1")
				break
			}

			// Dont allow digits to be less than 1 and more than 10
			if digits < 1 {
				screen.SetError(7, "Digits cannot be less than 1")
				break
			}

			// Dont allow digits to be less than 1 and more than 10
			if digits > 10 {
				screen.SetError(7, "Digits cannot be more than 10")
				break
			}

			// Instance of form
			formItems := screen.form.Items

			// Okay its time to add!
			token := tlockvault.Token{
				ID:               uuid.NewString(),
				Account:          formItems[0].FormItem.Value(),
				Issuer:           formItems[1].FormItem.Value(),
				Secret:           formItems[2].FormItem.Value(),
				Type:             toTokenType(formItems[3].FormItem.Value()),
				HashingAlgorithm: utils.ToHashFunction(formItems[4].FormItem.Value()),
				Period:           utils.Or(formItems[5].FormItem.Value(), 30),
				InitialCounter:   utils.Or(formItems[6].FormItem.Value(), 0),
				Digits:           utils.Or(formItems[7].FormItem.Value(), 6),
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

	// Update the height of the viewport
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))
	screen.viewport.Height = min(height, lipgloss.Height(screen.content))

	return screen, tea.Batch(cmds...)
}

// View
func (screen AddTokenScreen) View() string {
	return screen.viewport.View()
}

// Sets a custom error message for the given form item index
func (screen AddTokenScreen) SetError(itemIndex int, error string) {
	if item, ok := screen.form.Items[itemIndex].FormItem.(form.FormItemInputBox); ok {
		// Set the error
		item.ErrorMessage = &error

		// Update item
		screen.form.Items[itemIndex].FormItem = item
	}
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
