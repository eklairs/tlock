package config

import (
	"os"
	"path"

	_ "embed"

	bubblekey "github.com/charmbracelet/bubbles/key"
	"github.com/eklairs/tlock/tlock-internal/utils"
	"gopkg.in/yaml.v3"
)

// Returns the path to use for the user
func pathFor(user string) string {
	return path.Join(CONFIG_DIR, user, "keybindings.yaml")
}

// Default config
//
//go:embed keybindings.yaml
var DEFAULT_CONFIG_RAW []byte

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
	MoveUp Keybinding `yaml:"move_up"`

	// Move folder down
	MoveDown Keybinding `yaml:"move_down"`

	// Delete
	Delete Keybinding `yaml:"delete"`

	// Delete
	Copy Keybinding `yaml:"copy"`

	// Next token for HOTP
	NextHOTP Keybinding `yaml:"next_hotp"`

	// Next token for HOTP
	Move Keybinding `yaml:"move"`
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

// Load key bindings for a specific user
func LoadKeyBindings(user string) KeybindingsConfig {
	// Default key bindings
	default_keys := DefaultKeybindingsConfig()

	// Read file
	raw, err := os.ReadFile(pathFor(user))

	// If there is error, return the default keybindings
	if err != nil {
		// The file probably doesnt exist,
		// Lets write the default config
		WriteDefault(user)

		// Return default
		return default_keys
	}

	// Err
	if err := yaml.Unmarshal(raw, &default_keys); err != nil {
		return default_keys
	}

	return default_keys
}

// Writes the default keybindings configuration
func WriteDefault(user string) {
	// Open file
	if file, err := utils.EnsureExists(pathFor(user)); err == nil {
		file.Write(DEFAULT_CONFIG_RAW)
	}
}
