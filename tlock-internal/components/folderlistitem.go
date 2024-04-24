package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

// Active folder list item
func ActiveFolderListItem(name string, tokensCount int) string {
	ui := lipgloss.JoinVertical(
		lipgloss.Left,
		tlockstyles.Styles.Title.Render(name),
		tlockstyles.Styles.SubText.Render(fmt.Sprintf("%d tokens", tokensCount)),
	)

	return tlockstyles.Styles.FolderItemActive.Render(ui)
}

// Inactive folder list item
func InactiveFolderListItem(name string, tokensCount int) string {
	ui := lipgloss.JoinVertical(
		lipgloss.Left,
		tlockstyles.Styles.SubText.Render(name),
		tlockstyles.Styles.SubText.Render(fmt.Sprintf("%d tokens", tokensCount)),
	)

	return tlockstyles.Styles.FolderItemInactive.Render(ui)
}
