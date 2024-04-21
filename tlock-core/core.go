package tlockcore

import (
	"log"
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/kelindar/binary"

	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

// Users
var USERS_PATH = path.Join(xdg.DataHome, "tlock", "users")

// Represents a user
type UserSpec struct {
	// Username
	Username string

	// Path to its vault
	Vault string
}

// TLock core
type TLockCore struct {
	// Users
	Users []UserSpec
}

// Initializes a new instance of the core
func New() TLockCore {
	// Read users path
	raw, err := os.ReadFile(USERS_PATH)

	// Check for errors
	if err != nil {
		log.Printf("[users] Users map does not exist, returning empty users map")

		return TLockCore{}
	}

	users := TLockCore{}

	// Parse
	if err = binary.Unmarshal(raw, &users); err != nil {
		log.Fatalf("[users] Failed to parse users map, syntax error possibly: %v", err)
	}

	return users
}

// [PRIVATE] Writes the current users value to the file
func (users TLockCore) write() {
	// Serialize
	data, _ := binary.Marshal(users)

	// Create file
	file, err := os.Create(USERS_PATH)

	// Check for errors
	if err != nil {
		log.Fatalf("Failed to create users file")
	}

	// Write
	file.Write(data)
}

// Adds a new user
func (users *TLockCore) AddNewUser(username, password string) tlockvault.TLockVault {
	// Initialize new vault
	vault := tlockvault.Initialize(password)

	// Add to users
	users.Users = append(users.Users, UserSpec{Username: username, Vault: vault.Path})

	// Write
	users.write()

	// Return vault
	return vault
}
