package tlockvault

import (
	"errors"
	"os"
	"path"

	"github.com/kelindar/binary"
)

// Error represents the vault may be been moved or deleted
var ERR_VAULT_DELETED = errors.New("The vault does not exist, did you delete it?")

// Error represents that the password is invalid
var ERR_PASSWORD_INVALID = errors.New("Wrong password, please try again")

// Initializes a new instance of the vault at the given path
func Initialize(at, password string) (*Vault, error) {
	// Log if there was error while creating
	if err := os.MkdirAll(path.Dir(at), os.ModePerm); err != nil {
		return nil, err
	}

	// Initialize sender
	// Lets keep the size to 1 because we only want top data
	sender := make(chan []Folder, 1)

	// Initialize vault
	vault := Vault{
		path:     at,
		password: password,
		dataChan: sender,
	}

	// Start worker
	go vault.startFileWriterWorker(sender)

	// Write empty data
	vault.write()

	// Return
	return &vault, nil
}

// Loads a new vault instance
// Loads a vault instance from the given path
func Load(path, password string) (*Vault, error) {
	// Raw data
	var raw []byte
	var decrypted []byte

	// Any error
	var err error

	// Read encrypted bytes
	if raw, err = os.ReadFile(path); err != nil {
		return nil, ERR_VAULT_DELETED
	}

	// Empty data
	var data []Folder

	// Decrypt
	if decrypted, err = Decrypt(password, raw); err != nil {
		return nil, ERR_PASSWORD_INVALID
	}

	// Unmarshal binary serialized data
	if err := binary.Unmarshal(decrypted, &data); err != nil {
		return nil, ERR_PASSWORD_INVALID
	}

	// Sender channel
	// Lets keep the size to 1 because we only want top data
	sender := make(chan []Folder, 1)

	// Create vault instance and return
	vault := &Vault{
		path:     path,
		Folders:  data,
		password: password,
		dataChan: sender,
	}

	// Start worker
	go vault.startFileWriterWorker(sender)

	// Return
	return vault, nil
}
