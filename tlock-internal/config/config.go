package config

import (
	"encoding/json"
	"os"
	"path"

	"github.com/adrg/xdg"
)

// Default theme
var DEFAULT_THEME = "catppuccin"

// Path to the config file
var CONFIG_PATH = path.Join(xdg.ConfigHome, "tlock", "tlock.yaml")

// Represents the tlock config
type Config struct {
	// Current theme
	// Defaults to `Catppuccin`
	CurrentTheme string
}

// Returns the default config
func DefaultConfig() Config {
	return Config{
		CurrentTheme: DEFAULT_THEME,
	}
}

// Loads the config from the file
func GetConfig() Config {
	default_config := DefaultConfig()

	// Read raw
	config_raw, err := os.ReadFile(CONFIG_PATH)

	// If error, just return the default config
	if err != nil {
		return default_config
	}

	// Parse
	if err := json.Unmarshal(config_raw, &default_config); err != nil {
		return default_config
	}

	// Return
	return default_config
}