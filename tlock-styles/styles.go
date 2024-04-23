package tlockstyles

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	tlockthemes "github.com/eklairs/tlock/tlock-themes"
)

// Base styles that is common to all
type Styles struct {
	// Base
	Base lipgloss.Style

	// Title
	Title lipgloss.Style

	// Active list item
	ActiveItem lipgloss.Style

	// Inactive list item
	InactiveListItem lipgloss.Style

	// Center
	Center lipgloss.Style

	// Dimmed
	Dimmed lipgloss.Style

	// Input
	Input lipgloss.Style

	// Input placeholder
	InputPlaceholder lipgloss.Style

	// Error
	Error lipgloss.Style

	// Folder inactive
	FolderInactive lipgloss.Style

	// Folder active
	FolderActive lipgloss.Style

	// Title
	AccentTitle lipgloss.Style

	// Dimmed title
	DimmedTitle lipgloss.Style
}

// Initializes new instance of styles
func InitializeStyle(width int, theme tlockthemes.Theme) Styles {
	// Terminal option for setting background color to the theme's window bg
	renderer := lipgloss.NewRenderer(os.Stdout, termenv.WithProfile(termenv.TrueColor), func(o *termenv.Output) {
		// Convert it to rgba
		r, g, b, _ := theme.WindowBg.RGBA()

		// Set background color
		o.SetBackgroundColor(o.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b)))
	})

	// Base
	base := renderer.NewStyle().
		Width(width).
		Foreground(theme.WindowFg)

	// Base for list items
	listItem := base.Copy().Padding(1, 3).Foreground(theme.Dimmed)

	// Return
	return Styles{
		Base:             base,
		Center:           base.Copy().AlignHorizontal(lipgloss.Center),
		Dimmed:           base.Copy().Foreground(theme.Dimmed),
		Title:            base.Copy().Bold(true).Foreground(theme.Accent),
		ActiveItem:       listItem.Copy().Background(theme.WindowBgOver).Foreground(theme.Accent).Bold(true),
		InactiveListItem: listItem,
		Input:            base.Copy().Background(theme.WindowBgOver).Padding(1, 3),
		InputPlaceholder: base.Copy().UnsetWidth().Background(theme.WindowBgOver).Foreground(theme.Dimmed),
		Error:            base.Copy().Foreground(theme.Error),
		AccentTitle:      base.Copy().UnsetWidth().Background(theme.Accent).Padding(0, 1).Foreground(theme.WindowBg),
		DimmedTitle:      base.Copy().UnsetWidth().Background(theme.WindowBgOver).Padding(0, 1).Foreground(theme.WindowFg),
		FolderInactive:   base.Copy().Padding(1, 3),
		FolderActive: base.Copy().
			Background(theme.WindowBgOver).
			Padding(1, 2).
			BorderBackground(theme.WindowBgOver).
			BorderForeground(theme.Accent).
			Border(lipgloss.OuterHalfBlockBorder(), false, false, false, true),
	}
}
