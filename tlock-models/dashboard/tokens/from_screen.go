package tokens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	"github.com/pquerna/otp"

	"github.com/kbinani/screenshot"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
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
}

// ShortHelp()
func (k confirmScreenKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Continue, k.Retake}
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
}

type TokenFromScreen struct {
	// Help
	help help.Model

	// State
	state int

	// Token read from string
	token *string
}

// Initializes a new instance of fromScreen from screen
func InitializeTokenFromScreen() TokenFromScreen {
	return TokenFromScreen{
		help:  components.BuildHelp(),
		state: stateTake,
	}
}

// Init
func (screen TokenFromScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen TokenFromScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	var cmd tea.Cmd

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, fromScreenKeys.GoBack):
			if screen.state == stateTake {
				manager.PopScreen()
			}

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
			if screen.token != nil {
				cmd = func() tea.Msg {
					return AddTokenMessage{
						Token: *screen.token,
					}
				}
			}

			manager.PopScreen()
		}
	}

	return screen, cmd
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
			screen.help.View(fromScreenKeys),
		)
	case stateConfirm:
		items := []string{
			tlockstyles.Styles.Title.Render(fromScreenAsciiArt), "",
			tlockstyles.Styles.SubText.Render("Confirm addition of the token"), "",
		}

		if screen.token == nil {
			items = append(items, tlockstyles.Styles.Error.Render("Did not find any token!"))
		} else {
			key, err := otp.NewKeyFromURL(*screen.token)

			if err != nil {
				items = append(items, tlockstyles.Styles.Error.Render("Did not find any token!"))
			} else {
				items = append(items, fmt.Sprintf("%s %s", tlockstyles.Styles.SubText.Render("Found a token for "), tlockstyles.Styles.Title.Render(key.AccountName())))
			}
		}

		// Add help
		items = append(items, "", screen.help.View(confirmScreenKeys))

		return lipgloss.JoinVertical(
			lipgloss.Center,
			items...,
		)
	}

	return "Loading..."
}
