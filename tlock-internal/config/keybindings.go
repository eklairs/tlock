package config

import (
	"os"
	"path"

	"github.com/adrg/xdg"
	bubblekey "github.com/charmbracelet/bubbles/key"
	"gopkg.in/yaml.v3"
)

// Path to keybindings
var KEYBINDINGS_CONFIG = path.Join(xdg.ConfigHome, "tlock", "keybindings.yaml")

// Just a wrapper
type Keybinding struct {
	bubblekey.Binding
}

// Custom unmarshaller
func (key *Keybinding) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Keys
	var keys []string

	// Check for errors
	if err := unmarshal(&keys); err != nil {
		return err
	}

	// Parse
	key.Binding = bubblekey.NewBinding(bubblekey.WithKeys(keys...))

	// Return
	return nil
}

// Quick utility to create key binding
func new_key(keys ...string) Keybinding {
	return Keybinding{Binding: bubblekey.NewBinding(bubblekey.WithKeys(keys...))}
}

// Keybindings config
type KeybindingsConfig struct {
	// Folder keybindings
	Folder FolderKeyBinds `yaml:"folders"`

	// Token keybindings
	Tokens TokenKeyBinds `yaml:"tokens"`
}

// Folder keybinds
type FolderKeyBinds struct {
	// Add folder
	Add Keybinding `yaml:"add"`

	// Edit folder
	Edit Keybinding `yaml:"edit"`

	// Next
	Next Keybinding `yaml:"next"`

	// Previous
	Previous Keybinding `yaml:"previous"`

	// Move folder up
	MoveUp Keybinding `yaml:"move_up"`

	// Move folder down
	MoveDown Keybinding `yaml:"move_down"`

	// Delete
	Delete Keybinding `yaml:"delete"`
}

// Tokens keybinds
type TokenKeyBinds struct {
	// Add token
	Add Keybinding `yaml:"add"`

	// Add token from screen
	AddScreen Keybinding `yaml:"add_from_screen"`

	// Edit token
	Edit Keybinding `yaml:"edit"`

	// Next
	Next Keybinding `yaml:"next"`

	// Previous
	Previous Keybinding `yaml:"previous"`

	// Move folder up
	MoveUp Keybinding `yaml:"J"`

	// Move folder down
	MoveDown Keybinding `yaml:"K"`

	// Delete
	Delete Keybinding `yaml:"d"`

	// Delete
	Copy Keybinding `yaml:"c"`

	// Next token for HOTP
	NextHOTP Keybinding `yaml:"n"`

	// Next token for HOTP
	Move Keybinding `yaml:"m"`
}

// Returns the default keybindings
func DefaultKeybindingsConfig() KeybindingsConfig {
	return KeybindingsConfig{
		Folder: DefaultFolderKeyBinds(),
		Tokens: DefaultTokensKeyBinds(),
	}
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

// Load key bindings
func LoadKeyBindings() KeybindingsConfig {
	// Default key bindings
	default_keys := DefaultKeybindingsConfig()

	// Read file
	raw, err := os.ReadFile(KEYBINDINGS_CONFIG)

	// If there is error, return the default keybindings
	if err != nil {
		return default_keys
	}

	// Err
	if err := yaml.Unmarshal(raw, &default_keys); err != nil {
		return default_keys
	}

	return default_keys
}
