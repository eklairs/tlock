package tlockstyles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/context"
)

// Instance of styles
// Must be initialized on program's start
var Styles TLockStyles

// Themes used all over tlock
type TLockStyles struct {
    // Title
    Title lipgloss.Style

    // Sub text
    SubText lipgloss.Style
}

// Initializes the styles
func InitializeStyles(theme context.Theme) {
    base := lipgloss.NewStyle()

    // Initialize styles
    Styles = TLockStyles{
        Title: with(base).Foreground(theme.Accent).Bold(true),
        SubText: with(base).Foreground(theme.SubText),
    }
}

// Utility to copy a style
func with(style lipgloss.Style) lipgloss.Style {
	return style.Copy()
}
