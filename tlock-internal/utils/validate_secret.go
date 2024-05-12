package utils

import (
	"encoding/base32"
	"strings"
)

// Kanged from totp module
func ValidateSecret(secret string) bool {
	// Add pads
	secret = strings.TrimSpace(secret)
	if n := len(secret) % 8; n != 0 {
		secret = secret + strings.Repeat("=", 8-n)
	}

	// As noted in issue #24 Google has started producing base32 in lower case,
	// but the StdEncoding (and the RFC), expect a dictionary of only upper case letters.
	secret = strings.ToUpper(secret)

	// Check
	_, err := base32.StdEncoding.DecodeString(secret)

	// Return
	return err == nil
}
