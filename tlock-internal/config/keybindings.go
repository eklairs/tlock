package config

import (
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/bubbles/key"
	"gopkg.in/yaml.v3"
)

var KEYBINDINGS_CONFIG = path.Join(xdg.ConfigHome, "tlock", "keybindings.yaml")

// / Tries to fetch the key, runs it through the parsers or returns the default keys
func parse[T any](keys map[string]map[string][]string, key string, parser func(map[string][]string) T, default_value func() T) T {
	if value, ok := keys[key]; ok {
		return parser(value)
	}

	return default_value()
}

// Quick utility function for finding keys
func find(keymap map[string][]string, key string, default_value key.Binding) []string {
	if value, ok := keymap[key]; ok {
		return value
	}

	return default_value.Keys()
}

// / Quick utility to create key binding
func new_key(keys ...string) key.Binding {
	return key.NewBinding(key.WithKeys(keys...))
}

// Keybindings config
type KeybindingsConfig struct {
	// Folder keybindings
	Folder FolderKeyBinds

	// Token keybindings
	Tokens TokenKeyBinds
}

// / Returns the default keybindings
func DefaultKeybindingsConfig() KeybindingsConfig {
	return KeybindingsConfig{
		Folder: DefaultFolderKeyBinds(),
		Tokens: DefaultTokensKeyBinds(),
	}
}

// Folder keybinds
type FolderKeyBinds struct {
	// Add folder
	Add key.Binding

	// Edit folder
	Edit key.Binding

	// Next
	Next key.Binding

	// Previous
	Previous key.Binding

	// Move folder up
	MoveUp key.Binding

	// Move folder down
	MoveDown key.Binding

	// Delete
	Delete key.Binding
}

// Default folder keybindings
func DefaultFolderKeyBinds() FolderKeyBinds {
	return FolderKeyBinds{
		Add:      new_key("A"),
		Edit:     new_key("E"),
		Next:     new_key("tab"),
		Previous: new_key("shift+tab"),
		MoveUp:   new_key("ctrl+up"),
		MoveDown: new_key("ctrl+down"),
		Delete:   new_key("D"),
	}
}

// Parses a map to Folder keybindings
func ParseFolderKeyBinds(keys map[string][]string) FolderKeyBinds {
	// Get the default keys
	default_keys := DefaultFolderKeyBinds()

	return FolderKeyBinds{
		Add:      new_key(find(keys, "add", default_keys.Add)...),
		Edit:     new_key(find(keys, "edit", default_keys.Edit)...),
		Next:     new_key(find(keys, "next", default_keys.Next)...),
		Previous: new_key(find(keys, "previous", default_keys.Previous)...),
		MoveUp:   new_key(find(keys, "move_up", default_keys.MoveUp)...),
		MoveDown: new_key(find(keys, "move_down", default_keys.MoveDown)...),
		Delete:   new_key(find(keys, "delete", default_keys.Delete)...),
	}
}

// Tokens keybinds
type TokenKeyBinds struct {
	// Add token
	Add key.Binding

	// Add token from screen
	AddScreen key.Binding

	// Edit token
	Edit key.Binding

	// Next
	Next key.Binding

	// Previous
	Previous key.Binding

	// Move folder up
	MoveUp key.Binding

	// Move folder down
	MoveDown key.Binding

	// Delete
	Delete key.Binding

	// Delete
	Copy key.Binding

	// Next token for HOTP
	NextHOTP key.Binding

	// Next token for HOTP
	Move key.Binding
}

// Default tokens keybindings
func DefaultTokensKeyBinds() TokenKeyBinds {
	return TokenKeyBinds{
		Add:       new_key("a"),
		Edit:      new_key("e"),
		Next:      new_key("j"),
		Previous:  new_key("k"),
		MoveUp:    new_key("J"),
		MoveDown:  new_key("K"),
		Delete:    new_key("d"),
		AddScreen: new_key("s"),
		Copy:      new_key("c"),
		Move:      new_key("m"),
		NextHOTP:  new_key("n"),
	}
}

// Parses a map to Folder keybindings
func ParseTokensKeyBinds(keys map[string][]string) TokenKeyBinds {
	// Get the default keys
	default_keys := DefaultTokensKeyBinds()

	return TokenKeyBinds{
		Add:       new_key(find(keys, "add", default_keys.Add)...),
		Edit:      new_key(find(keys, "edit", default_keys.Edit)...),
		Next:      new_key(find(keys, "next", default_keys.Next)...),
		Previous:  new_key(find(keys, "previous", default_keys.Previous)...),
		MoveUp:    new_key(find(keys, "move_up", default_keys.MoveUp)...),
		MoveDown:  new_key(find(keys, "move_down", default_keys.MoveDown)...),
		Delete:    new_key(find(keys, "delete", default_keys.Delete)...),
		AddScreen: new_key(find(keys, "add_from_screen", default_keys.AddScreen)...),
		Copy:      new_key(find(keys, "copy", default_keys.Copy)...),
		NextHOTP:  new_key(find(keys, "next_hotp", default_keys.NextHOTP)...),
		Move:      new_key(find(keys, "move", default_keys.Move)...),
	}
}

// Load key bindings
func LoadKeyBindings() KeybindingsConfig {
	// Read file
	raw, err := os.ReadFile(KEYBINDINGS_CONFIG)

	// If there is error, return the default keybindings
	if err != nil {
		return DefaultKeybindingsConfig()
	}

	// Parse
	var mapped_keys map[string]map[string][]string

	// Err
	if err := yaml.Unmarshal(raw, &mapped_keys); err != nil {
		return DefaultKeybindingsConfig()
	}

	// Keybindings
	return KeybindingsConfig{
		Folder: parse(mapped_keys, "folders", ParseFolderKeyBinds, DefaultFolderKeyBinds),
		Tokens: parse(mapped_keys, "tokens", ParseTokensKeyBinds, DefaultTokensKeyBinds),
	}
}
