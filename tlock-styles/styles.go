package tlockstyles

import "github.com/eklairs/tlock/tlock-internal/context"

// Instance of styles
// Must be initialized on program's start
var Styles TLockStyles

// Themes used all over tlock
type TLockStyles struct {

}

// Initializes the styles
func InitializeStyles(theme context.Theme) {
    Styles = TLockStyles{

    }
}
