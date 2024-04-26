package tokens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/tlock-messages"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"github.com/kbinani/screenshot"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/pquerna/otp"
)

var fromScreenAsciiArt = `
█▀ █▀▀ █▀█ █▀▀ █▀▀ █▄ █
▄█ █▄▄ █▀▄ ██▄ ██▄ █ ▀█`

const (
	stateTake = iota
	stateConfirm
)

// From screen key map
type fromScreenKeyMap struct {
	GoBack key.Binding
	Start  key.Binding
}

// ShortHelp()
func (k fromScreenKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.GoBack, k.Start}
}

// FullHelp()
func (k fromScreenKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.GoBack},
		{k.Start},
	}
}

// Keys
var fromScreenKeys = fromScreenKeyMap{
	GoBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
	Start: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "start"),
	),
}

// Confirm from screen keys
type confirmScreenKeyMap struct {
	Continue key.Binding
	Retake   key.Binding
	Escape   key.Binding
}

// ShortHelp()
func (k confirmScreenKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Escape, k.Continue, k.Retake}
}

// FullHelp()
func (k confirmScreenKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

// Keys
var confirmScreenKeys = confirmScreenKeyMap{
	Continue: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "continue"),
	),
	Retake: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "retake"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
}

type TokenFromScreen struct {
	// State
	state int

	// Vault
	vault *tlockvault.Vault

	// Token read from string
	token *string

	// Folder
	folder tlockvault.Folder
}

// Initializes a new instance of fromScreen from screen
func InitializeTokenFromScreen(vault *tlockvault.Vault, folder tlockvault.Folder) TokenFromScreen {
	return TokenFromScreen{
		state:  stateTake,
		vault:  vault,
		folder: folder,
	}
}

// Init
func (screen TokenFromScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen TokenFromScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, fromScreenKeys.GoBack):
			manager.PopScreen()

		case key.Matches(msgType, fromScreenKeys.Start) && screen.state == stateTake:
			if image, err := screenshot.CaptureRect(screenshot.GetDisplayBounds(0)); err == nil {
				if bmp, err := gozxing.NewBinaryBitmapFromImage(image); err == nil {
					qrReader := qrcode.NewQRCodeReader()

					if result, err := qrReader.Decode(bmp, nil); err == nil {
						uri := result.String()

						screen.token = &uri
					}
				}
			}

			screen.state = stateConfirm

		case key.Matches(msgType, confirmScreenKeys.Retake):
			if screen.state == stateConfirm {
				screen.state = stateTake
			}

		case key.Matches(msgType, confirmScreenKeys.Continue) && screen.state == stateConfirm:
			// Add the token
			if screen.token != nil {
				// Add token
				screen.vault.AddToken(screen.folder.ID, *screen.token)

				// Require refresh of folders and tokens list
				cmds = append(
					cmds,
					func() tea.Msg { return tlockmessages.RefreshFoldersMsg{} },
					func() tea.Msg { return tlockmessages.RefreshTokensMsg{} },
				)
			}

			manager.PopScreen()
		}
	}

	return screen, tea.Batch(cmds...)
}

// View
func (screen TokenFromScreen) View() string {
	switch screen.state {
	case stateTake:
		return lipgloss.JoinVertical(
			lipgloss.Center,
			tlockstyles.Styles.Title.Render(fromScreenAsciiArt), "",
			tlockstyles.Styles.SubText.Render("Place your QRCode window in the same screen as tlock"), "",
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				tlockstyles.Styles.MockScreen.Render("TLock"), "   ",
				tlockstyles.Styles.MockScreen.Render("QRCode Window"),
			), "",
			tlockstyles.Help.View(fromScreenKeys),
		)

	case stateConfirm:
		items := []string{
			tlockstyles.Styles.Title.Render(fromScreenAsciiArt), "",
			tlockstyles.Styles.SubText.Render("Confirm addition of the token"), "",
		}

		// If the token is null, show the message
		if screen.token == nil {
			items = append(items, tlockstyles.Styles.Error.Render("Did not find any token!"))
		} else {
			// Try to parse the otp value
			key, err := otp.NewKeyFromURL(*screen.token)

			// If there was error while finding
			if err != nil {
				// Show the error
				items = append(items, tlockstyles.Styles.Error.Render("Did not find any token!"))

				// Reset token screen before its not what we were looking for
				screen.token = nil
			} else {
				// Find the account name
				accountName := key.AccountName()

				if accountName == "" {
					accountName = "<not found>"
				}

				// Show to user
				items = append(items, fmt.Sprintf("%s %s", tlockstyles.Styles.SubText.Render("Found a token for"), tlockstyles.Styles.Title.Render(accountName)))
			}
		}

		// Add help
		items = append(items, "", tlockstyles.Help.View(confirmScreenKeys))

		return lipgloss.JoinVertical(
			lipgloss.Center,
			items...,
		)
	}

	return "Loading..."
}
