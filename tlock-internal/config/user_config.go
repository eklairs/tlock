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
	return path.Join(CONFIG_DIR, user, "config.yaml")
}

// Default config
//
//go:embed default_config.yaml
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
type UserConfiguration struct {
	// Whether to enable icons
	EnableIcons bool `yaml:"enable_icons"`

	// Folder keybindings
	Folder FolderKeyBinds `yaml:"folders_keybindings"`

	// Token keybindings
	Tokens TokenKeyBinds `yaml:"tokens_keybindings"`
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
func DefaultUserConfiguration() UserConfiguration {
	return UserConfiguration{
		EnableIcons: false,
		Folder:      DefaultFolderKeyBinds(),
		Tokens:      DefaultTokensKeyBinds(),
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

// Load configuration for a specific user
func LoadUserConfig(user string) UserConfiguration {
	// Default key bindings
	default_config := DefaultUserConfiguration()

	// Parse the file if it is read
	if raw, err := os.ReadFile(pathFor(user)); err == nil {
		yaml.Unmarshal(raw, &default_config)
	}

	// Return
	return default_config
}

// Writes the default keybindings configuration
func WriteDefault(user string) {
	// Open file
	if file, err := utils.EnsureExists(pathFor(user)); err == nil {
		file.Write(DEFAULT_CONFIG_RAW)
	}
}
