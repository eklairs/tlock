package models

import (
	tlockcore "github.com/eklairs/tlock/tlock-core"
	tlockthemes "github.com/eklairs/tlock/tlock-themes"
)

// TLock Context
type Context struct {
    // Current theme
    Theme tlockthemes.Theme

    // Core
    Core tlockcore.TLockCore
}

// Initializes a new instance of the context
func InitializeContext() Context {
    return Context{
        Theme: tlockthemes.CATPPUCCIN_THEME,
        Core: tlockcore.New(),
    }
}

