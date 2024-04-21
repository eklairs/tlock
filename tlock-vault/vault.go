package tlockvault

import (
	"log"
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/google/uuid"
	"github.com/kelindar/binary"
)

// Removes an index from a slice
func remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

// Represents a folder
type FolderSpec struct {
	// Name
	Name string

	// Tokens uris
	Uris []string
}

// Data inside of the vault
type TLockVaultData struct {
	Folders []FolderSpec
}

// Vault
type TLockVault struct {
	// Path to the vault file
	Path string

	// Password to encrypt with
	password string

	// Data
	Data TLockVaultData
}

// Initializes a new vault at the given path
func Initialize(password string) TLockVault {
	// New uuid
	id := uuid.New()
	dir := path.Join(xdg.DataHome, "tlock", "root", id.String())

	// Make root dir
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create user's root dir: %v", err)
	}

	// Initialize vault
	vault := TLockVault{
		Path:     path.Join(dir, "vault"),
		password: password,
	}

	// Write empty data
	vault.write()

	// Return
	return vault
}

// Loads a vault instance from the given path
func Load(path, password string) (*TLockVault, error) {
	// Read encrypted bytes
	raw, err := os.ReadFile(path)

	// No errors, pl0x
	if err != nil {
		return nil, err
	}

	// Empty data
	data := TLockVaultData{}

	// Decrypt
	decrypted, err := Decrypt(password, raw)
	if err != nil {
		return nil, err
	}

	// Unmarshal binary serialized data
	if err := binary.Unmarshal(decrypted, &data); err != nil {
		return nil, err
	}

	// Create vault instance and return
	return &TLockVault{
		Data:     data,
		Path:     path,
		password: password,
	}, nil
}

// Writes the current value of the vault
func (vault TLockVault) write() {
	// Serialize
	serialized, _ := binary.Marshal(vault.Data)

	// Encrypt
	encrypted := Encrypt(vault.password, serialized)

	// Create file
	f, _ := os.Create(vault.Path)

	// Write
	if _, err := f.Write(encrypted); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
}
