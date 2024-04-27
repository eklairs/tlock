package config

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/rs/zerolog/log"
)

// Default theme
var DEFAULT_THEME = "Dracula"

// Path to the config file
var CONFIG_PATH = path.Join(xdg.ConfigHome, "tlock", "tlock.json")

// Represents theme config
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
		// Log
		log.Debug().Msg("[config] No config file found, returning the default config")

		return default_config
	}

	// Parse
	if err := json.Unmarshal(config_raw, &default_config); err != nil {
		// Log
		log.Error().Err(err).Msg("[config] Failed to parse config, syntax error possibly?")

		return default_config
	}

	// Return
	return default_config
}

// Writes the config
func (config Config) Write() {
	// Make directory
	os.MkdirAll(filepath.Dir(CONFIG_PATH), os.ModePerm)

	// Marshal
	data, _ := json.Marshal(config)

	// Write
	file, err := os.Create(CONFIG_PATH)

	// If no error, write to file
	if err == nil {
		file.Write(data)
	} else {
		// Log
		log.Error().Err(err).Msg("[config] Failed to write to config")
	}
}
