package tlockvault

import (
	"log"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

// Size of the salt
var SALT_SIZE = 32

// Derives a new key from the password to use it for cryptographic purposes using argon2id
// You can pass salt which will be used, or let the function generate one for you
func GenerateKey(password string, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = make([]byte, SALT_SIZE)

		if _, err := rand.Read(salt); err != nil {
			log.Fatalf("[key_generator] Failed to generate random salt: %v", err)
		}
	}

	return argon2.Key([]byte(password), salt, 3, 32*1024, 4, 32), salt
}

// Encrypts the given piece of byte array
func Encrypt(password string, data []byte) []byte {
	// Ciphers and blocks
	var blockCipher cipher.Block
	var gcm cipher.AEAD

	// Error
	var err error

	key, salt := GenerateKey(password, nil)

	// Initialize AES
	if blockCipher, err = aes.NewCipher(key); err != nil {
		log.Fatalf("[encrypt] Failed to create a new AES cipher: %v", err)
	}

	// Initialize GCM
	if gcm, err = cipher.NewGCM(blockCipher); err != nil {
		log.Fatalf("[encrypt] Failed to create a new GCM block: %v", err)
	}

	// Encrypt
	nonce := make([]byte, gcm.NonceSize())

	encryptedText := gcm.Seal(nonce, nonce, data, nil)
	encryptedText = append(encryptedText, salt...)

	return encryptedText
}

// Decrypts the given piece of encrypted byte array
// It returns an error if decryption fails, because of the invalid key
func Decrypt(password string, data []byte) ([]byte, error) {
	// Ciphers and blocks
	var blockCipher cipher.Block
	var gcm cipher.AEAD

	var err error

	// Extract salt and data
	salt, data := data[len(data)-SALT_SIZE:], data[:len(data)-SALT_SIZE]

	key, salt := GenerateKey(password, salt)

	// Initialize AES
	if blockCipher, err = aes.NewCipher(key); err != nil {
		log.Fatalf("[decrypt] Failed to create a new AES cipher: %v", err)
	}

	// Initialize GCM
	if gcm, err = cipher.NewGCM(blockCipher); err != nil {
		log.Fatalf("[decrypt] Failed to create a new GCM block: %v", err)
	}

	// Get nonce
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	// Decrypt
	decryptedText, err := gcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return nil, err
	}

	return decryptedText, nil
}
