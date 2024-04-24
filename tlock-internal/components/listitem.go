package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

// List item active
func ListItemActive(width int, title, suffix string) string {
	space_width := width - lipgloss.Width(title) - lipgloss.Width(suffix)

	ui := lipgloss.JoinHorizontal(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(title),
		tlockstyles.Styles.SubAltBg.Render(strings.Repeat(" ", space_width)),
		tlockstyles.Styles.SubAltBg.Render(tlockstyles.Styles.Title.Render(suffix)),
	)

	return tlockstyles.Styles.ListItemActive.Render(ui)
}

// List item active
func ListItemInactive(width int, title, suffix string) string {
	space_width := width - lipgloss.Width(title) - lipgloss.Width(suffix)

	ui := lipgloss.JoinHorizontal(
		lipgloss.Center,
		tlockstyles.Styles.SubText.Render(title),
		strings.Repeat(" ", space_width),
		tlockstyles.Styles.SubText.Render(suffix),
	)

	return tlockstyles.Styles.ListItemInactive.Render(ui)
}
