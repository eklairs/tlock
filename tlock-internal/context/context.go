package context

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	tlockcore "github.com/eklairs/tlock/tlock-core"
	"github.com/eklairs/tlock/tlock-internal/config"
	tlockvendor "github.com/eklairs/tlock/tlock-vendor"
	"golang.design/x/clipboard"
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
	Config config.Config

	// Core
	Core tlockcore.TLockCore

	// If the clipboard is available
	ClipboardAvailability bool

	// Keybindings
	Keybindings config.KeybindingsConfig
}

// Initializes a new instance of the context
// It is recommended to call this at a one place and then pass around the context
func InitializeContext() Context {
	// Init clipboard
	err := clipboard.Init()

	// Themes
	var themes []Theme
	json.Unmarshal(tlockvendor.ThemesJSON, &themes)

	// Parse icons
	var icons struct {
		Icons map[string]Icon
	}

	json.Unmarshal(tlockvendor.IconsJSON, &icons)

	// Return
	return Context{
		Themes:                themes,
		Icons:                 icons.Icons,
		Core:                  tlockcore.New(),
		Config:                config.GetConfig(),
		ClipboardAvailability: err == nil,
		Keybindings:           config.LoadKeyBindings(),
	}
}

// Finds the index of the theme
func (context Context) findTheme(name string) int {
	return slices.IndexFunc(context.Themes, func(theme Theme) bool { return strings.ToLower(theme.Name) == strings.ToLower(name) })
}

// Returns the current theme spec
func (context Context) GetCurrentTheme() Theme {
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
