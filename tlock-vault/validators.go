package tlockvault

import (
	"errors"
	"strings"

	"github.com/eklairs/tlock/tlock-internal/utils"
)

// Error representing that the folder name is empty
var ERR_FOLDER_EMPTY = errors.New("Folder name cannot be empty")

// Error representing that the folder with that name already exists
var ERR_FOLDER_EXISTS = errors.New("Folder with that name already exists")

// Error representing that the token secret is empty
var ERR_TOKEN_EMPTY = errors.New("Secret value cannt be empty")

// Error representing that the secret is invalid
var ERR_TOKEN_INVALID = errors.New("Secret is invalid, are you sure it is typed correctly?")

// Error representing that the secret already exists
var ERR_TOKEN_EXISTS = errors.New("Token with that secret already exists")

// Validates if the folder name is fit to be used
func (vault Vault) validateFolderName(name string) (string, error) {
	// Sanitize by trimming off the spaces
	name = strings.TrimSpace(name)

	// Check if the folder name is not empty
	if name == "" {
		return name, ERR_FOLDER_EMPTY
	}

	// Check if the folder already exists
	if vault.folderExists(name) {
		return name, ERR_FOLDER_EXISTS
	}

	// Return
	return name, nil
}

// Validates if the token is fit to be used
// It is checked on the basis of the fact that it can be used to generate a secret
// And no other token with the same secret exist
func (vault Vault) ValidateToken(secret string) (string, error) {
	// Sanitize by trimming off the spaces
	secret = strings.TrimSpace(secret)

	// Check if the folder name is not empty
	if secret == "" {
		return secret, ERR_TOKEN_EMPTY
	}

	// Check if the folder already exists
	if vault.tokenExists(secret) {
		return secret, ERR_TOKEN_EXISTS
	}

	// Try to generate token
	if !utils.ValidateSecret(secret) {
		return secret, ERR_TOKEN_INVALID
	}

	// Return
	return secret, nil
}
