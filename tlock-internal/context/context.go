package context

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	tlockcore "github.com/eklairs/tlock/tlock-core"
	"github.com/eklairs/tlock/tlock-internal/config"
	tlockvendor "github.com/eklairs/tlock/tlock-vendor"
)

type Icon struct {
	Unicode string
	Hex     string
}

// Represents a theme
type Theme struct {
	// Name
	Name string

	// Background
	Background lipgloss.Color

	// Background over
	BackgroundOver lipgloss.Color

	// Sub text
	SubText lipgloss.Color

	// Accent
	Accent lipgloss.Color

	// Foreground
	Foreground lipgloss.Color

	// Error
	Error lipgloss.Color
}

// Represents a context
type Context struct {
	// All the themes available
	// Fetched from vendor
	Themes []Theme

	// Icons!
	Icons map[string]Icon

	// Config
	TLockConfig config.TLockConfig

	// Core
	Core *tlockcore.TLockCore

	// User configuration
	Config config.UserConfiguration
}

// Initializes a new instance of the context
// It is recommended to call this at a one place and then pass around the context
func InitializeContext() Context {
	// Themes
	var themes []Theme
	json.Unmarshal(tlockvendor.ThemesJSON, &themes)

	// Parse icons
	var icons struct {
		Icons map[string]Icon
	}

	json.Unmarshal(tlockvendor.IconsJSON, &icons)

	// Initialize core
	core, _ := tlockcore.New()

	// Return
	return Context{
		Themes:      themes,
		Icons:       icons.Icons,
		Core:        core,
		Config:      config.DefaultUserConfiguration(),
		TLockConfig: config.GetTLockConfig(),
	}
}

// Finds the index of the theme
func (context Context) findTheme(name string) int {
	return slices.IndexFunc(context.Themes, func(theme Theme) bool { return strings.ToLower(theme.Name) == strings.ToLower(name) })
}

// Returns the current theme spec
func (context Context) GetCurrentTheme() Theme {
	// Get theme index
	theme_index := context.findTheme(context.TLockConfig.CurrentTheme)

	// If not found, then use the default theme
	if theme_index == -1 {
		theme_index = context.findTheme(config.DEFAULT_THEME)
	}

	// Return
	return context.Themes[theme_index]
}

// Sets the theme
func (context *Context) SetTheme(theme string) {
	// Update current theme
	context.TLockConfig.CurrentTheme = theme

	// Write
	context.TLockConfig.Write()
}
