package tokens

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/messages"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-internal/utils"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	tlockstyles "github.com/eklairs/tlock/tlock/styles"
	"github.com/pquerna/otp"
)

var MeterV2 = spinner.Spinner{
	Frames: []string{
		"▱▱▱▱",
		"▰▱▱▱",
		"▰▰▱▱",
		"▰▰▰▱",
		"▰▰▰▰",
		"▱▰▰▰",
		"▱▱▰▰",
		"▱▱▱▰",
	},
	FPS: time.Second / 8, //nolint:gomnd
}

// Channel to send the token read from the screen
var dataFromScreenChan = make(chan *dataFromScreen)

type dataFromScreen struct {
	// The otp
	Uri *string

	// Any error from validators
	Err error
}

// Message stating that a data has been recved
type dataRecievedMsg struct {
	data *dataFromScreen
}

var fromScreenAsciiArt = `
█▀ █▀▀ █▀█ █▀▀ █▀▀ █▄ █
▄█ █▄▄ █▀▄ ██▄ ██▄ █ ▀█`

const (
	stateTake = iota
	stateGathering
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

func pollDataFetched() tea.Cmd {
	return func() tea.Msg {
		return dataRecievedMsg{
			data: <-dataFromScreenChan,
		}
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

	// Spinner
	spinner spinner.Model

	// Token read from string
	token *dataFromScreen

	// Folder
	folder tlockvault.Folder

	// Status bar message to send
	statusBarMessage string
}

// Initializes a new instance of fromScreen from screen
func InitializeTokenFromScreen(vault *tlockvault.Vault, folder tlockvault.Folder) TokenFromScreen {
	// Initialize spinner
	s := spinner.New()
	s.Spinner = MeterV2
	s.Style = tlockstyles.Styles.Title

	// Return
	return TokenFromScreen{
		state:   stateTake,
		vault:   vault,
		spinner: s,
		folder:  folder,
	}
}

// Init
func (screen TokenFromScreen) Init() tea.Cmd {
	return pollDataFetched()
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
			screen.state = stateGathering

			// Start spinner
			cmds = append(cmds, screen.spinner.Tick)

			go func() {
				// Read data from the screen
				data, err := utils.ReadTokenFromScreen()

				// Prepare data
				dataScreen := dataFromScreen{
					Uri: data,
					Err: err,
				}

				// Try to parse
				if data != nil {
					if key, err := otp.NewKeyFromURL(*data); err == nil {
						// Run validator
						_, err := screen.vault.ValidateToken(key.Secret())

						// Update error message
						if dataScreen.Err == nil {
							dataScreen.Err = err
						}
					}
				}

				// Send
				dataFromScreenChan <- &dataScreen
			}()

		case key.Matches(msgType, confirmScreenKeys.Retake):
			if screen.state == stateConfirm {
				// Set state
				screen.state = stateTake

				// Restart poll
				cmds = append(cmds, pollDataFetched())
			}

		case key.Matches(msgType, confirmScreenKeys.Continue) && screen.state == stateConfirm:
			// Add the token
			if screen.token != nil && screen.token.Err == nil {
				// Add token
				// We can ignore validations because we have already pre-checked it
				screen.vault.AddToken(screen.folder.Name, *screen.token.Uri)

				// Require refresh of folders and tokens list
				cmds = append(
					cmds,
					func() tea.Msg { return tlockmessages.RefreshFoldersMsg{} },
					func() tea.Msg { return tlockmessages.RefreshTokensMsg{} },
					func() tea.Msg { return components.StatusBarMsg{Message: screen.statusBarMessage} },
				)
			}

			manager.PopScreen()
		}

	case dataRecievedMsg:
		screen.token = msgType.data
		screen.state = stateConfirm
	}

	if screen.state == stateGathering {
		var cmd tea.Cmd
		screen.spinner, cmd = screen.spinner.Update(msg)
		cmds = append(cmds, cmd)
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

	case stateGathering:
		return screen.spinner.View()

	case stateConfirm:
		items := []string{
			tlockstyles.Styles.Title.Render(fromScreenAsciiArt), "",
			tlockstyles.Styles.SubText.Render("Confirm addition of the token"), "",
		}

		// If the token is null, show the message
		if screen.token != nil {
			// Try to parse the otp value
			if screen.token.Err == nil {
				key, err := otp.NewKeyFromURL(*screen.token.Uri)

				// If there was error while finding
				if err != nil {
					// Show the error
					items = append(items, tlockstyles.Styles.Error.Render("Did not find any token!"))

					// Reset token screen before its not what we were looking for
					screen.token = nil
				} else {
					// Find the account name
					accountName := key.AccountName()
					screen.statusBarMessage = fmt.Sprintf("Successfully added token for %s from screen", accountName)

					if accountName == "" {
						accountName = "<no account name>"
						screen.statusBarMessage = fmt.Sprintf("Successfully added token from screen (no account name)")
					}

					// Show to user
					items = append(items, fmt.Sprintf("%s %s", tlockstyles.Styles.SubText.Render("Found a token for"), tlockstyles.Styles.Title.Render(accountName)))
				}
			} else {
				items = append(items, lipgloss.JoinHorizontal(
					lipgloss.Center,
					tlockstyles.Styles.Error.Render(screen.token.Err.Error()),
				))
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
