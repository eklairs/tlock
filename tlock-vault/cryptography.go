package tlockvault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

// Size of the salt
var SALT_SIZE = 32

// Derives a new key from the password to use it for cryptographic purposes using argon2id
// You can pass salt which will be used, or let the function generate one for you
// It returns (key, salt, error)
func GenerateKey(password string, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, SALT_SIZE)

		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	return argon2.Key([]byte(password), salt, 3, 32*1024, 4, 32), salt, nil
}

// Encrypts the given piece of byte array
func Encrypt(password string, data []byte) ([]byte, error) {
	// Ciphers and blocks
	var blockCipher cipher.Block
	var gcm cipher.AEAD
	var err error

	// Key
	var key []byte
	var salt []byte

	// Generate key
	if key, salt, err = GenerateKey(password, nil); err != nil {
		return nil, err
	}

	// Initialize AES
	if blockCipher, err = aes.NewCipher(key); err != nil {
		return nil, err
	}

	// Initialize GCM
	if gcm, err = cipher.NewGCM(blockCipher); err != nil {
		return nil, err
	}

	// Encrypt
	nonce := make([]byte, gcm.NonceSize())

	// Encrypt
	encryptedText := gcm.Seal(nonce, nonce, data, nil)
	encryptedText = append(encryptedText, salt...)

	// Return
	return encryptedText, nil
}

// Decrypts the given piece of encrypted byte array
// It returns an error if decryption fails, because of the invalid key
func Decrypt(password string, data []byte) ([]byte, error) {
	// Ciphers and blocks
	var blockCipher cipher.Block
	var gcm cipher.AEAD
	var decryptedText []byte
	var err error

	// Key
	var key []byte

	// Extract salt and data
	salt, data := data[len(data)-SALT_SIZE:], data[:len(data)-SALT_SIZE]

	// Generate key
	if key, _, err = GenerateKey(password, salt); err != nil {
		return nil, err
	}

	// Initialize AES
	if blockCipher, err = aes.NewCipher(key); err != nil {
		return nil, err
	}

	// Initialize GCM
	if gcm, err = cipher.NewGCM(blockCipher); err != nil {
		return nil, err
	}

	// Get nonce
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	// Decrypt
	if decryptedText, err = gcm.Open(nil, nonce, ciphertext, nil); err != nil {
		return nil, err
	}

	return decryptedText, nil
}
