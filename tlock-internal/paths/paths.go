package paths

import (
	"path"

	"github.com/adrg/xdg"
)

// Base for data directory
var DATA_BASE = path.Join(xdg.DataHome, "tlock")

// Base for config directory
var CONFIG_BASE = path.Join(xdg.ConfigHome, "tlock")

// Directory that contains vaults
var VAULT_DIR = path.Join(DATA_BASE, "vaults")

// Path to the list of users
var USERS = path.Join(DATA_BASE, "users.bin")

// Path to tlock internal config file
var TLOCK_CONFIG = path.Join(CONFIG_BASE, "config_internal_ignore.bin")

// Returns the path to the given user's config
func UserConfigFor(username string) string {
	return path.Join(CONFIG_BASE, username, "config.yaml")
}
