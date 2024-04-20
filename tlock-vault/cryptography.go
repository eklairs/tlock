package tlockvault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

var SALT_SIZE = 32

// Derives a new key from the password to use it for cryptographic purposes
// You can pass salt which will be used, or let the function generate one for you
func GenerateKey(password string, salt []byte) ([]byte, []byte) {
    if salt == nil {
        salt = make([]byte, SALT_SIZE)

        if _, err := rand.Read(salt); err != nil {
            panic("Failed to generate random number")
        }
    }

    return argon2.Key([]byte(password), salt, 3, 32 * 1024, 4, 32), salt
}

// Encrypts the given piece of byte array
func Encrypt(password string, data []byte) []byte {
    key, salt := GenerateKey(password, nil);

    blockCipher, _ := aes.NewCipher(key)
    gcm, _ := cipher.NewGCM(blockCipher)
    nonce := make([]byte, gcm.NonceSize())

    encryptedText := gcm.Seal(nonce, nonce, data, nil)
    encryptedText = append(encryptedText, salt...)

    return encryptedText
}

// Decrypts the given piece of encrypted byte array
func Decrypt(password string, data []byte) ([]byte, error) {
    salt, data := data[len(data) - SALT_SIZE:], data[:len(data) - SALT_SIZE]

    key, _ := GenerateKey(password, salt)

    blockCipher, _ := aes.NewCipher(key)
    gcm, _ := cipher.NewGCM(blockCipher)

    nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

    decryptedText, err := gcm.Open(nil, nonce, ciphertext, nil)

    if err != nil {
        return nil, err
    }

    return decryptedText, nil
}

