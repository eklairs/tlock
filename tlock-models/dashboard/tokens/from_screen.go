package tokens

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/buildhelp"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	"github.com/kbinani/screenshot"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

var TOKEN_FROM_SCREEN_WIDTH = 65

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

type TokenFromScreen struct {
	// Styles
	styles tlockstyles.Styles

	// Help
	help help.Model

	// Context
	context context.Context
}

// Initializes a new instance of fromScreen from screen
func InitializeTokenFromScreen(context context.Context) TokenFromScreen {
	// Initialize styles
	styles := tlockstyles.InitializeStyle(TOKEN_FROM_SCREEN_WIDTH, context.Theme)

	// IBuild help
	help := buildhelp.BuildHelp(styles)

	return TokenFromScreen{
		styles:  styles,
		help:    help,
		context: context,
	}
}

// Init
func (model TokenFromScreen) Init() tea.Cmd {
	return nil
}

// Update
func (model TokenFromScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	var cmd tea.Cmd

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, fromScreenKeys.GoBack):
			manager.PopScreen()

		case key.Matches(msgType, fromScreenKeys.Start):
			if image, err := screenshot.CaptureRect(screenshot.GetDisplayBounds(0)); err == nil {
				if bmp, err := gozxing.NewBinaryBitmapFromImage(image); err == nil {
					qrReader := qrcode.NewQRCodeReader()

					if result, err := qrReader.Decode(bmp, nil); err == nil {
						cmd = func() tea.Msg {
							return AddTokenMsg{
								URI: result.String(),
							}
						}
					}
				}
			}

			manager.PopScreen()
		}
	}

	return model, cmd
}

// View
func (model TokenFromScreen) View() string {
	mockScreenStyle := model.styles.Base.Copy().
		Background(model.context.Theme.WindowBgOver).
		Align(lipgloss.Center, lipgloss.Center).
		Width(27).
		Height(9)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		model.styles.Center.Render(model.styles.Dimmed.Render("Place your QRCode window in the same screen as tlock")), "",
		model.styles.Center.Render(lipgloss.JoinHorizontal(
			lipgloss.Left,
			mockScreenStyle.Render("TLock"),
			"    ",
			mockScreenStyle.Render("QRCode Window"),
		)), "",
		model.styles.Center.Render(model.help.View(fromScreenKeys)),
	)
}
