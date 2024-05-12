package tlockcore

import (
	"os"
	"path"
	"slices"

	"github.com/eklairs/tlock/tlock-internal/paths"
	"github.com/eklairs/tlock/tlock-internal/utils"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"github.com/kelindar/binary"
)

// Core utilities
type TLockCore struct {
	// Users
	Users []User
}

// [PRIVATE] Writes the current users value to the file
func (users TLockCore) write() {
	// Serialize
	data, _ := binary.Marshal(users)

	// Create file
	if file, err := utils.EnsureExists(paths.USERS); err == nil {
		file.Write(data)
	}
}

// Initializes a new instance of the core
func New() (*TLockCore, error) {
	// Raw
	var raw []byte
	var err error

	// Read users path
	if raw, err = os.ReadFile(paths.USERS); err != nil {
		return &TLockCore{}, err
	}

	// Parsed users
	var users []User

	// Parse
	if err = binary.Unmarshal(raw, &users); err != nil {
		return &TLockCore{}, nil
	}

	return &TLockCore{Users: users}, nil
}

// Adds a new user
func (core *TLockCore) AddNewUser(username, password string) (*tlockvault.Vault, error) {
	var err error

	// Vault
	var vault *tlockvault.Vault

	// Run validations
	if username, err = core.validateUsername(username); err != nil {
		return nil, err
	}

	// New user
	user := User(username)

	// Initialize new vault
	if vault, err = tlockvault.Initialize(user.Vault(), password); err != nil {
		return nil, err
	}

	// Add to users
	core.Users = append(core.Users, user)

	// Write
	core.write()

	// Return vault
	return vault, nil
}

// Renames a user
func (core *TLockCore) RenameUser(old, new string) error {
	var err error

	// Run validations
	if new, err = core.validateUsername(new); err != nil {
		return err
	}

	// Find
	if index := core.Find(old); index != -1 {
		// Original vault of the user
		original_vault := path.Dir(core.Users[index].Vault())

		// Rename vault folder
		if err := os.Rename(original_vault, path.Join(path.Dir(original_vault), new)); err != nil {
			return err
		}

		// Rename if the user is found
		core.Users[index] = User(new)

		// Write
		core.write()
	}

	// No errors
	return nil
}

// Deletes a user
func (core *TLockCore) DeleteUser(username string) {
	if index := core.Find(username); index != -1 {
		// Remove its vault
		os.RemoveAll(path.Dir(User(username).Vault()))

		// Remove its config
		os.RemoveAll(path.Dir(paths.UserConfigFor(username)))

		// Remove user
		core.Users = utils.Remove(core.Users, index)

		// Write
		core.write()
	}
}

// Returns if the user with the given name already exists
func (core TLockCore) Exists(username string) bool {
	return core.Find(username) != -1
}

// Returns the index of the user
func (core TLockCore) Find(username string) int {
	return slices.Index(core.Users, User(username))
}
