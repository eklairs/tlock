package components

import (
	"fmt"
	"strings"

	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

// Active folder list item
func ActiveFolderListItem(name string, tokensCount int) string {
	items := []string{
		tlockstyles.Styles.Title.Render(name),
		tlockstyles.Styles.SubText.Render(fmt.Sprintf("%d tokens", tokensCount)),
	}

	return tlockstyles.Styles.FolderItemActive.Render(strings.Join(items, "\n"))
}

// Inactive folder list item
func InactiveFolderListItem(name string, tokensCount int) string {
	items := []string{
		tlockstyles.Styles.SubText.Render(name),
		tlockstyles.Styles.SubText.Render(fmt.Sprintf("%d tokens", tokensCount)),
	}

	return tlockstyles.Styles.FolderItemInactive.Render(strings.Join(items, "\n"))
}
