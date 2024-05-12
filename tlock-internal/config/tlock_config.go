package config

import (
	"os"

	"github.com/eklairs/tlock/tlock-internal/paths"
	"github.com/eklairs/tlock/tlock-internal/utils"
	"github.com/kelindar/binary"
)

// Default theme
var DEFAULT_THEME = "Catppuccin"

// TLock config is the config which is overriden by tlock itself
type TLockConfig struct {
	// Current theme
	// Defaults to `Catppuccin`
	CurrentTheme string `yaml:"current_theme"`
}

// Returns the default config
func DefaultTLockConfig() TLockConfig {
	return TLockConfig{
		CurrentTheme: DEFAULT_THEME,
	}
}

// Loads the config from the file
func GetTLockConfig() TLockConfig {
	default_config := DefaultTLockConfig()

	// Read raw
	if config_raw, err := os.ReadFile(paths.TLOCK_CONFIG); err == nil {
		binary.Unmarshal(config_raw, &default_config)
	}

	// Return
	return default_config
}

// Writes the config
func (config TLockConfig) Write() {
	// Marshal
	data, _ := binary.Marshal(config)

	// Create file
	if file, err := utils.EnsureExists(paths.TLOCK_CONFIG); err == nil {
		file.Write(data)
	}
}
