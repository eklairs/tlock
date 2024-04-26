package tlockstyles

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/context"
)

// Help
var Help help.Model

// Instance of styles
// Must be initialized on program's start
var Styles TLockStyles

// Themes used all over tlock
type TLockStyles struct {
	// Title
	Title lipgloss.Style

	// Sub text
	SubText lipgloss.Style

	// Style for items over SubAltBg
	SubAltBg lipgloss.Style

	// Style for active list items
	ListItemActive lipgloss.Style

	// Style for inactive list item
	ListItemInactive lipgloss.Style

	// Style for input
	Input lipgloss.Style

	// Error
	Error lipgloss.Style

	// Style for placeholder
	Placeholder lipgloss.Style

	// Folder active list item
	FolderItemActive lipgloss.Style

	// Folder inactive list item
	FolderItemInactive lipgloss.Style

	// Tilte Bar
	AccentBgItem lipgloss.Style

	// Mock screen
	MockScreen lipgloss.Style

	// Background Over
	BackgroundOver lipgloss.Style
}

// Initializes the styles
func InitializeStyles(theme context.Theme) {
	// Base
	base := lipgloss.NewStyle()

	// Base for padded items
	paddedItem := with(base).Padding(1, 3)

	// Initialize styles
	Styles = TLockStyles{
		Title:              with(base).Foreground(theme.Accent).Bold(true),
		SubText:            with(base).Foreground(theme.SubText),
		SubAltBg:           with(base).Background(theme.BackgroundOver),
		Error:              with(base).Foreground(theme.Error).Bold(true),
		Input:              with(paddedItem).Width(65).Background(theme.BackgroundOver),
		Placeholder:        with(base).Background(theme.BackgroundOver).Foreground(theme.SubText),
		ListItemActive:     with(paddedItem).Background(theme.BackgroundOver),
		AccentBgItem:       with(base).Padding(0, 1).Background(theme.Accent).Foreground(theme.Background),
		BackgroundOver:     with(base).Background(theme.BackgroundOver),
		FolderItemInactive: with(paddedItem),
		MockScreen:         with(base).Background(theme.BackgroundOver).Align(lipgloss.Center, lipgloss.Center).Width(27).Height(9),
		FolderItemActive: with(paddedItem).
			Padding(1, 2).
			Background(theme.BackgroundOver).
			Border(lipgloss.OuterHalfBlockBorder(), false, false, false, true).
			BorderBackground(theme.BackgroundOver).
			BorderForeground(theme.Accent),
		ListItemInactive: with(paddedItem),
	}

	// Initialize help
	Help = help.New()

	// Comply help menu styles to themes
	Help.Styles.ShortKey = Styles.Title
	Help.Styles.ShortDesc = Styles.SubText
	Help.Styles.ShortSeparator = Styles.SubText
}

// Utility to copy a style
func with(style lipgloss.Style) lipgloss.Style {
	return style.Copy()
}
