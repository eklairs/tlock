package tlockvault

import (
	"slices"

	"github.com/eklairs/tlock/tlock-internal/utils"
	"github.com/pquerna/otp"
)

// Converts `totp` or `hotp` to TokenType
func toType(type_ string) TokenType {
	if type_ == "hotp" {
		return TokenTypeHOTP
	}

	return TokenTypeTOTP
}

// Adds a new token to the given folder from token URI
func (vault *Vault) AddToken(folder string, uri string) error {
	// Key
	var key *otp.Key
	var err error

	// Generate key
	if key, err = otp.NewKeyFromURL(uri); err == nil {
		// Generate token
		token := Token{
			Type:             toType(key.Type()),
			Issuer:           key.Issuer(),
			Account:          key.AccountName(),
			Secret:           key.Secret(),
			InitialCounter:   0,
			Period:           int(key.Period()),
			Digits:           key.Digits().Length(),
			HashingAlgorithm: key.Algorithm(),
			UsageCounter:     0,
		}

		// Add
		return vault.AddTokenFromToken(folder, token)
	}

	// Return
	return err
}

// Adds a new token to the given folder
func (vault *Vault) AddTokenFromToken(folder string, token Token) error {
	var err error

	if token.Secret, err = vault.ValidateToken(token.Secret); err == nil {
		// Find folder and if it exists, add
		if index := vault.findFolder(folder); index != -1 {
			vault.Folders[index].Tokens = append(vault.Folders[index].Tokens, token)
		}

		// Write
		vault.write()
	}

	// Return
	return err
}

// Replace a token in the given folder
func (vault *Vault) ReplaceToken(fromFolder string, token, newToken Token) error {
	var err error

	if newToken.Secret, err = vault.ValidateToken(newToken.Secret); token.Secret == newToken.Secret || err == nil {
		// Get folder index
		if index := vault.findFolder(fromFolder); index != -1 {
			// Replace
			vault.Folders[index].Tokens[vault.findToken(index, token.Secret)] = newToken

			// Write
			vault.write()

			// Ok!
			return nil
		}
	}

	return err
}

// Deletes a token in the given folder
func (vault *Vault) DeleteToken(folder string, token Token) {
	// Find the folder
	if folder, token := vault.locateToken(folder, token); folder != -1 && token != -1 {
		vault.Folders[folder].Tokens = utils.Remove(vault.Folders[folder].Tokens, token)
	}

	// Write
	vault.write()
}

// Move a token to the given folder
func (vault *Vault) MoveToken(token Token, fromFolder, toFolder string) {
	// Remove from existing
	vault.DeleteToken(fromFolder, token)

	// Add to the new one
	vault.AddTokenFromToken(toFolder, token)
}

// Move a token to the given folder
func (vault *Vault) IncreaseCounter(folder string, token Token) {
	// Find the folder
	if folder, token := vault.locateToken(folder, token); folder != -1 && token != -1 {
		vault.Folders[folder].Tokens[token].UsageCounter++
	}

	// Write
	vault.write()
}

// Moves the token down
func (vault *Vault) MoveTokenDown(folder string, token Token) bool {
	// Find
	if folder, token := vault.locateToken(folder, token); folder != -1 && token != -1 {
		// If it is already at the bottom, skip
		if token == len(vault.Folders[folder].Tokens)-1 {
			return false
		}

		// Swapppp
		vault.Folders[folder].Tokens = utils.Swap(vault.Folders[folder].Tokens, token, token+1)

		// Wrap
		vault.write()

		// Return
		return true
	}

	return false
}

// Moves the token up
func (vault *Vault) MoveTokenUp(folder string, token Token) bool {
	// Find
	if folder, token := vault.locateToken(folder, token); folder != -1 && token != -1 {
		// If it is already at the bottom, skip
		if token == 0 {
			return false
		}

		// Swapppp
		vault.Folders[folder].Tokens = utils.Swap(vault.Folders[folder].Tokens, token, token-1)

		// Wrap
		vault.write()

		// Return
		return true
	}

	return false
}

// Find a token index by its secret
func (vault *Vault) findToken(folder int, secret string) int {
	return slices.IndexFunc(vault.Folders[folder].Tokens, func(token Token) bool { return token.Secret == secret })
}

// Checks if the token exists with a secret
func (vault *Vault) tokenExists(secret string) bool {
	for i := 0; i < len(vault.Folders); i++ {
		if vault.findToken(i, secret) != -1 {
			return true
		}
	}

	return false
}

// Finds the index of folder as well token
func (vault *Vault) locateToken(folder string, token Token) (int, int) {
	// Find folder index
	folderIndex := vault.findFolder(folder)

	// Return
	return folderIndex, vault.findToken(folderIndex, token.Secret)
}
