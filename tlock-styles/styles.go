package tlockstyles

import (
	"github.com/charmbracelet/lipgloss"
)

// Instance of the styles
// Must call `InitializeStyles()` to initialize these styles
var Styles TLockStyles

// Represents a theme
type Theme struct {
	// Name
	Name string

	// Background color
	Background lipgloss.Color

	// Main Color
	Accent lipgloss.Color

	// Sub color
	Sub lipgloss.Color

	// Sub Alt color
	SubAlt lipgloss.Color

	// Text color
	Text lipgloss.Color

	// Error
	Error lipgloss.Color
}

// Styles used by tlock
type TLockStyles struct {
	// Title
	Title lipgloss.Style

	// SubText
	SubText lipgloss.Style

	// Style for active list items
	ListItemActive lipgloss.Style

	// Style for inactive list item
	ListItemInactive lipgloss.Style
}

// Initializes the styles
func InitializeStyles(theme Theme) {
	// Base that every style must copy from
	base := lipgloss.NewStyle()

	// Base for padded items
	paddedItem := with(base).Padding(1, 3)

	// Initialize styles
	Styles = TLockStyles{
		Title:            with(base).Foreground(theme.Accent),
		SubText:          with(base).Foreground(theme.Sub),
		ListItemActive:   with(paddedItem).Background(theme.SubAlt),
		ListItemInactive: with(paddedItem),
	}
}

// Utility to copy a style
func with(style lipgloss.Style) lipgloss.Style {
	return style.Copy()
}
