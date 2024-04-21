package buildhelp

import (
	"github.com/charmbracelet/bubbles/help"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

func BuildHelp(styles tlockstyles.Styles) help.Model {
    // Help menu
	help := help.New()

	// Comply help menu styles to themes
	// Unsetting width because we want the menus to let occupy the needed space
	help.Styles.ShortKey = styles.Title.Copy().UnsetWidth()
	help.Styles.ShortDesc = styles.Dimmed.Copy().UnsetWidth()
	help.Styles.ShortSeparator = styles.Dimmed.Copy().UnsetWidth()

    return help
}
