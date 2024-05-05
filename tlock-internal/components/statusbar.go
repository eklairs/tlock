package components

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/constants"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	"golang.org/x/term"
)

type StatusBarMsg struct {
	Message      string
	ErrorMessage bool
}

// Message to remove status bar message
type RemoveStatusBarMsg struct{}

type StatusBar struct {
	// Message to show
	Message string

	// Is the meessage a error message
	ErrorMessage bool

	// Current user
	CurrentUser string
}

func NewStatusBar(currentUser string) StatusBar {
	return StatusBar{
		Message:     fmt.Sprintf("Welcome, %s!", currentUser),
		CurrentUser: currentUser,
	}
}

func (bar *StatusBar) Update(msg tea.Msg) tea.Cmd {
	switch msgType := msg.(type) {
	case StatusBarMsg:
		bar.Message = msgType.Message
		bar.ErrorMessage = msgType.ErrorMessage

		return func() tea.Msg {
			// Let the message be alive for 4 second
			time.Sleep(time.Second * 4)

			// Remove
			return RemoveStatusBarMsg{}
		}

	case RemoveStatusBarMsg:
		bar.Message = ""
	}

	return nil
}

func (bar *StatusBar) View() string {
	// Get width
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))

	items := make([]string, 5)

	// Add app name
	items[0] = tlockstyles.Styles.AccentBgItem.Render("TLOCK")

	// Add version
	items[1] = tlockstyles.Styles.OverlayItem.Render(constants.VERSION)

	// Current date, maybe?
	items[3] = tlockstyles.Styles.OverlayItem.Render(time.Now().Format("2 January, 2006"))

	// Current logged in user
	items[4] = tlockstyles.Styles.AccentBgItem.Render(bar.CurrentUser)

	for _, item := range items {
		width -= lipgloss.Width(item)
	}

	// Render message
	messageStyle := tlockstyles.Styles.SubAltBg.Copy().Width(width)

	if bar.ErrorMessage {
		messageStyle = tlockstyles.Styles.Error.Copy().Inherit(tlockstyles.Styles.SubAltBg).Width(width)
	}

	items[2] = messageStyle.Render(bar.Message)

	return lipgloss.JoinHorizontal(lipgloss.Left, items...)
}
