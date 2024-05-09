package tlockcore

import (
	"errors"
	"strings"
)

// Error representing that the username is empty
var ERR_USERNAME_EMPTY = errors.New("Username cannot be empty")

// Error representing that the username is empty
var ERR_USERNAME_EXISTS = errors.New("User with that name already exists")

// Validates if the username is okay to be used
// Validation is done on the basis of its emptiness and if there are other users with the same name
func (core TLockCore) validateUsername(username string) (string, error) {
	// Strip off
	username = strings.TrimSpace(username)

	// Check if it is empty
	if username == "" {
		return username, ERR_USERNAME_EMPTY
	}

	// Check if it already exists
	if core.Exists(username) {
		return username, ERR_USERNAME_EXISTS
	}

	// No errors
	return username, nil
}

