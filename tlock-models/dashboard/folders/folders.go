package folders

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/eklairs/tlock/tlock-internal/context"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

// Folders
type Folders struct {
	// Context
	context context.Context

	// Vault
	vault *tlockvault.TLockVault

	// List view
	listview list.Model
}

// Initializes a new instance of folders
func InitializeFolders(vault *tlockvault.TLockVault, context context.Context) Folders {
	return Folders{
		vault:   vault,
		context: context,
	}
}
