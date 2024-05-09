package components

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/constants"
	tlockstyles "github.com/eklairs/tlock/tlock/styles"
	"golang.org/x/term"
)

type StatusBarMsg struct {
	Message      string
	ErrorMessage bool
}

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

	// Sub width to fill in the gap
	for _, item := range items {
		width -= lipgloss.Width(item)
	}

	// Render message
	items[2] = tlockstyles.Styles.SubAltBg.Copy().Width(width).Render()

	// Join them all
	ui := lipgloss.JoinHorizontal(lipgloss.Left, items...)

	// Message renderable
	render_fn := tlockstyles.Styles.Base.Render

	// If it is an error message
	if bar.ErrorMessage {
		render_fn = tlockstyles.Styles.Error.Render
	}

	// Return
	return lipgloss.JoinVertical(lipgloss.Left, ui, render_fn("› "+bar.Message))
}
