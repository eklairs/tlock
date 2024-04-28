package tlockcore

import (
	"errors"
	"os"
	"path"
	"slices"

	"github.com/adrg/xdg"
	"github.com/kelindar/binary"
	"github.com/rs/zerolog/log"

	"github.com/eklairs/tlock/tlock-internal/utils"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

// Users
var USERS_PATH = path.Join(xdg.DataHome, "tlock", "users")

// Represents a user
type User struct {
	// Username
	Username string

	// Path to its vault
	Vault string
}

// TLock core
type TLockCore struct {
	// Users
	Users []User
}

// Initializes a new instance of the core
func New() TLockCore {
	// Read users path
	raw, err := os.ReadFile(USERS_PATH)

	// Check for errors
	if err != nil {
		// Log
		log.Debug().Msg("[users] Users map does not exist, returning empty users map")

		// Return empty core
		return TLockCore{}
	}

	var users TLockCore

	// Parse
	if err = binary.Unmarshal(raw, &users); err != nil {
		log.Fatal().Err(err).Msg("[users] Failed to parse users map, syntax error possibly")
	}

	return users
}

// Adds a new user
func (users *TLockCore) AddNewUser(username, password string) (*tlockvault.Vault, error) {
	if users.exists(username) {
		return nil, errors.New("User already exists")
	}

	// Initialize new vault
	vault := tlockvault.Initialize(password)

	// Add to users
	users.Users = append(users.Users, User{Username: username, Vault: vault.Path})

	// Write
	users.write()

	// Return vault
	return &vault, nil
}

// [PRIVATE] Writes the current users value to the file
func (users TLockCore) write() {
	// Serialize
	data, _ := binary.Marshal(users)

	// Create file
	file, err := utils.EnsureExists(USERS_PATH)

	// Check for errors
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create users file")
	}

	// Write
	file.Write(data)
}

// Checks if a user with the name exists
func (users TLockCore) exists(username string) bool {
	return slices.IndexFunc(users.Users, func(user User) bool { return user.Username == username }) != -1
}
