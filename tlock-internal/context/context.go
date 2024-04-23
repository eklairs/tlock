package context

import (
	tlockcore "github.com/eklairs/tlock/tlock-core"
	tlockthemes "github.com/eklairs/tlock/tlock-themes"
	"golang.design/x/clipboard"
)

// TLock Context
type Context struct {
	// Current theme
	Theme tlockthemes.Theme

	// Core
	Core tlockcore.TLockCore

    // Is clipboard available
    ClipboardAvailability bool
}

// Initializes a new instance of the context
func InitializeContext() Context {
    err := clipboard.Init()

	return Context{
		Theme: tlockthemes.CATPPUCCIN_THEME,
		Core:  tlockcore.New(),
        ClipboardAvailability: err == nil,
	}
}
