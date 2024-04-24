package components

import (
	"github.com/charmbracelet/bubbles/help"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

func BuildHelp() help.Model {
	// Help menu
	help := help.New()

	// Comply help menu styles to themes
	// Unsetting width because we want the menus to let occupy the needed space
	help.Styles.ShortKey = tlockstyles.Styles.Title
	help.Styles.ShortDesc = tlockstyles.Styles.SubText
	help.Styles.ShortSeparator = tlockstyles.Styles.SubText

	return help
}
