package context

import (
	"encoding/json"
	"slices"
	"strings"

	tlockcore "github.com/eklairs/tlock/tlock-core"
	"github.com/eklairs/tlock/tlock-internal/config"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	tlockvendor "github.com/eklairs/tlock/tlock-vendor"
	"golang.design/x/clipboard"
)

// Context contains all the themes, config etc for the app
type Context struct {
	// All the themes available
	Themes []tlockstyles.Theme

	// Config
	Config config.Config

	// Core
	Core tlockcore.TLockCore

	// If the clipboard is available
	ClipboardAvailability bool
}

// Initializes a new instance of the context
// It is recommended to call this at a one place and then pass around the context
func InitializeContext() Context {
	err := clipboard.Init()

	themes := make([]tlockstyles.Theme, 0)

	// Parse themes
	json.Unmarshal(tlockvendor.ThemesJSON, &themes)

	// Return
	return Context{
		Themes:                themes,
		Core:                  tlockcore.New(),
		Config:                config.GetConfig(),
		ClipboardAvailability: err == nil,
	}
}

// Finds the index of the theme
func (context Context) findTheme(name string) int {
	return slices.IndexFunc(context.Themes, func(theme tlockstyles.Theme) bool { return strings.ToLower(theme.Name) == strings.ToLower(name) })
}

// Returns the current theme spec
func (context Context) GetCurrentTheme() tlockstyles.Theme {
	// Get theme index
	theme_index := context.findTheme(context.Config.CurrentTheme)

	// If not found, then use the default theme
	if theme_index == -1 {
		theme_index = context.findTheme(config.DEFAULT_THEME)
	}

	// Return
	return context.Themes[theme_index]
}

// Sets the theme
func (context *Context) SetTheme(theme string) {
    context.Config.CurrentTheme = theme

    context.Config.Write()
}
