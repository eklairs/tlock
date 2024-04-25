package tlockstyles

import (
	"math"
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
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
	// Base
	Base lipgloss.Style

	// Title
	Title lipgloss.Style

	// SubText
	SubText lipgloss.Style

	// Style for active list items
	ListItemActive lipgloss.Style

	// Style for inactive list item
	ListItemInactive lipgloss.Style

	// Style for items over SubAltBg
	SubAltBg lipgloss.Style

	// Style for placeholder
	Placeholder lipgloss.Style

	// Style for input
	Input lipgloss.Style

	// Error
	Error lipgloss.Style

	// Folder active list item
	FolderItemActive lipgloss.Style

	// Folder inactive list item
	FolderItemInactive lipgloss.Style

	// Tilte Bar
	TitleBar lipgloss.Style
}

// Initializes the styles
func InitializeStyles(theme Theme) {
	// Use 1/4 of the screen afor folder width
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))
	foldersWidth := int(math.Floor((1.0 / 5.0) * float64(width)))

	// Base that every style must copy from
	base := lipgloss.NewStyle().Foreground(theme.Text)

	// Base for padded items
	paddedItem := with(base).Padding(1, 3)

	// Initialize styles
	Styles = TLockStyles{
		Base:               with(base),
		Title:              with(base).Foreground(theme.Accent).Bold(true),
		SubText:            with(base).Foreground(theme.Sub),
		SubAltBg:           with(base).Background(theme.SubAlt),
		Placeholder:        with(base).Background(theme.SubAlt).Foreground(theme.Sub),
		Error:              with(base).Foreground(theme.Error).Bold(true),
		TitleBar:           with(base).Padding(0, 1).Background(theme.Accent).Foreground(theme.Background),
		Input:              with(paddedItem).Background(theme.SubAlt),
		ListItemActive:     with(paddedItem).Background(theme.SubAlt),
		ListItemInactive:   with(paddedItem),
		FolderItemInactive: with(paddedItem).Width(foldersWidth),
		FolderItemActive: with(paddedItem).
			Padding(1, 2).
			Background(theme.SubAlt).
			Width(foldersWidth).
			Border(lipgloss.OuterHalfBlockBorder(), false, false, false, true).
			BorderBackground(theme.SubAlt).
			BorderForeground(theme.Accent),
	}
}

// Utility to copy a style
func with(style lipgloss.Style) lipgloss.Style {
	return style.Copy()
}
