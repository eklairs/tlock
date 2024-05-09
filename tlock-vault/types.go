package tlockvault

import "github.com/pquerna/otp"

// Token types
const (
	TokenTypeTOTP = iota
	TokenTypeHOTP
)

// Token Type
type TokenType int

// Token
type Token struct {
	// Type
	Type TokenType

	// Issuer name
	Issuer string

	// Account name
	Account string

	// Secret
	Secret string

	// Initial counter [only in case of HOTP based tokens]
	InitialCounter int

	// Period [only in case of TOTP based tokens]
	Period int

	// Digits
	Digits int

	// Hasing function
	HashingAlgorithm otp.Algorithm

	// Usage counter [only in case of HOTP based tokens]
	UsageCounter int
}

// Folder
type Folder struct {
	// Name of the folder
	Name string

	// Tokens
	Tokens []Token
}
