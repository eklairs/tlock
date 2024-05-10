package tokens

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/form"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/messages"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-internal/utils"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	tlockstyles "github.com/eklairs/tlock/tlock/styles"
	"github.com/pquerna/otp"
	"golang.org/x/term"
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

var editTokenAscii = `
█▀▀ █▀▄ █ ▀█▀
██▄ █▄▀ █  █`

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

// Initializes a new screen of EditTokenScreen
func InitializeEditTokenScreen(folder tlockvault.Folder, token tlockvault.Token, vault *tlockvault.Vault) EditTokenScreen {
	tokenTypes := []string{"TOTP", "HOTP"}
	hashFunctions := []string{"SHA1", "SHA256", "SHA512", "MD5"}

	typeToString := func(tokenType tlockvault.TokenType) string {
		if tokenType == tlockvault.TokenTypeTOTP {
			return "TOTP"
		}

		if tokenType == tlockvault.TokenTypeHOTP {
			return "HOTP"
		}

		return ""
	}

	hashFunctionToString := func(hashFn otp.Algorithm) string {
		if hashFn == otp.AlgorithmMD5 {
			return "MD5"
		}

		if hashFn == otp.AlgorithmSHA1 {
			return "SHA1"
		}

		if hashFn == otp.AlgorithmSHA256 {
			return "SHA256"
		}

		if hashFn == otp.AlgorithmSHA512 {
			return "SHA512"
		}

		return ""
	}

	items := []form.FormItem{
		form.FormItemInputBox{
			Title:       "Account Name",
			Description: "Name of the account, like John Doe",
			Input:       value(components.InitializeInputBox("Account name goes here..."), token.Account),
		},
		form.FormItemInputBox{
			Title:       "Issuer Name",
			Description: "Name of the issuer, like GitHub",
			Input:       value(components.InitializeInputBox("Issuer name goes here..."), token.Issuer),
		},
		form.FormItemInputBox{
			Title:       "Secret",
			Description: "The secret provided by the issuer",
			Input:       value(components.InitializeInputBox("The secret goes here..."), token.Secret),
		},
		form.FormItemOptionBox{
			Title:         "Type",
			Description:   "Time or counter based token",
			Values:        tokenTypes,
			SelectedIndex: slices.Index(tokenTypes, typeToString(token.Type)),
		},
		form.FormItemOptionBox{
			Title:         "Hash function",
			Description:   "Hash function for the token",
			Values:        hashFunctions,
			SelectedIndex: slices.Index(hashFunctions, hashFunctionToString(token.HashingAlgorithm)),
		},
		form.FormItemInputBox{
			Title:       "Period",
			Description: "Time to refresh the token",
			Input:       utils.ValidatorIntegerNo0(value(components.InitializeInputBoxCustomWidth("Time in seconds...", 24), fmt.Sprintf("%d", token.Period))),
		},
		form.FormItemInputBox{
			Title:       "Usage counter",
			Description: "Usage counter of HOTP token",
			Input:       utils.ValidatorIntegerNo0(value(components.InitializeInputBoxCustomWidth("Usage counter", 24), fmt.Sprintf("%d", token.UsageCounter))),
		},
		form.FormItemInputBox{
			Title:       "Digits",
			Description: "Number of digits",
			Input:       utils.ValidatorIntegerNo0(value(components.InitializeInputBoxCustomWidth("Number of digits goes here...", 24), fmt.Sprintf("%d", token.Digits))),
		},
	}

	// Disable the HOTP initially
	form := form.New(items)
	form.Items[6].Enabled = false

	// Get term size
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))

	content := GenerateEditUI(form)

	// Initialize viewport
	viewport := utils.DisableViewportKeys(viewport.New(85, min(height, lipgloss.Height(content))))
	viewport.SetContent(content)

	return EditTokenScreen{
		form:     form,
		token:    token,
		vault:    vault,
		folder:   folder,
		content:  content,
		viewport: viewport,
	}
}

// Init
func (screen EditTokenScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen EditTokenScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, editTokenKeys.Enter):
			for index, item := range screen.form.Items {
				if input, ok := item.FormItem.(form.FormItemInputBox); ok {
					input.ErrorMessage = nil
					screen.form.Items[index].FormItem = input
				}
			}

			period := utils.Or(screen.form.Items[5].FormItem.Value(), screen.token.Period)
			digits := utils.Or(screen.form.Items[7].FormItem.Value(), screen.token.Digits)

			// Dont allow period to be zero
			if period < 1 {
				screen.SetError(5, errors.New("Period cannot be less than 1"))
				break
			}

			// Dont allow digits to be less than 1 and more than 10
			if digits < 1 {
				screen.SetError(7, errors.New("Digits cannot be less than 1"))
				break
			}

			// Dont allow digits to be less than 1 and more than 10
			if digits > 10 {
				screen.SetError(7, errors.New("Digits cannot be more than 10"))
				break
			}

			// Instance of form
			formItems := screen.form.Items

			// Okay its time to edit!
			token := tlockvault.Token{
				Account:          formItems[0].FormItem.Value(),
				Issuer:           formItems[1].FormItem.Value(),
				Secret:           formItems[2].FormItem.Value(),
				Type:             toTokenType(formItems[3].FormItem.Value()),
				HashingAlgorithm: utils.ToHashFunction(formItems[4].FormItem.Value()),
				Period:           period,
				InitialCounter:   screen.token.InitialCounter,
				Digits:           utils.Or(formItems[7].FormItem.Value(), screen.token.Digits),
				UsageCounter:     utils.Or(formItems[6].FormItem.Value(), screen.token.UsageCounter),
			}

			// Replace
            if err := screen.vault.ReplaceToken(screen.folder.Name, screen.token, token); err != nil {
                screen.SetError(2, err)
                break
            }

			accountName := formItems[0].FormItem.Value()

			statusBarMessage := fmt.Sprintf("Successfully edited token for %s", accountName)

			if accountName == "" {
				accountName = "<no account name>"
				statusBarMessage = fmt.Sprintf("Successfully edited token (no account name)")
			}

			// Require refresh of folders and tokens list
			cmds = append(
				cmds,
				func() tea.Msg { return tlockmessages.RefreshFoldersMsg{} },
				func() tea.Msg { return tlockmessages.RefreshTokensMsg{} },
				func() tea.Msg { return components.StatusBarMsg{Message: statusBarMessage} },
			)

			// Pop
			manager.PopScreen()

		case key.Matches(msgType, editTokenKeys.GoBack):
			manager.PopScreen()
		}

	case tea.WindowSizeMsg:
		screen.viewport.Height = min(msgType.Height, lipgloss.Height(screen.content))
	}

	// Update the form
	screen.form.Update(msg)

	// Get size
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))

	// Update viewport
	screen.viewport, _ = screen.viewport.Update(msg)

	// Generate UI
	screen.content = GenerateEditUI(screen.form)

	// Set viewport content
	screen.viewport.Height = min(lipgloss.Height(screen.content), height)
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
	screen.viewport.Height = min(height, lipgloss.Height(screen.content))

	return screen, tea.Batch(cmds...)
}

// View
func (screen EditTokenScreen) View() string {
	return screen.viewport.View()
}

// Generates the UI
func GenerateEditUI(form form.Form) string {
	// Items
	items := []string{
		tlockstyles.Styles.Title.Render(editTokenAscii), "",
		tlockstyles.Styles.SubText.Render("Edit a new token [secret is required, rest are optional]"), "",
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

	// Edit the help menu
	items = append(items, inputGroup, "", tlockstyles.Help.View(editTokenKeys))

	// Return!
	return lipgloss.JoinVertical(
		lipgloss.Center,
		items...,
	)
}

// Sets a custom error message for the given form item index
func (screen EditTokenScreen) SetError(itemIndex int, error error) {
	if item, ok := screen.form.Items[itemIndex].FormItem.(form.FormItemInputBox); ok {
		// Set the error
		item.ErrorMessage = &error

		// Update item
		screen.form.Items[itemIndex].FormItem = item
	}
}
