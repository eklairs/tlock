package tlockcore

import (
	"path"

	"github.com/eklairs/tlock/tlock-internal/paths"
)

// Represents a user
type User string

// Returns the string representation
func (user User) S() string {
	return string(user)
}

// Returns the path to its vault
func (user User) Vault() string {
	return path.Join(paths.VAULT_DIR, user.S(), "vault.bin")
}

