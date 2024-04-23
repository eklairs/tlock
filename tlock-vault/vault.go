package tlockvault

import (
	"log"
	"os"
	"path"
	"slices"

	"github.com/adrg/xdg"
	"github.com/google/uuid"
	"github.com/kelindar/binary"
)

// Removes an index from a slice
func remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

type TokenURI struct {
    // URI
    URI string

    // Usage count
    // Only in case of HOTP based tokens
    UsageCounter int
}

// Represents a folder
type FolderSpec struct {
	// Name
	Name string

	// Tokens uris
	Uris []TokenURI
}

// Data inside of the vault
type TLockVaultData struct {
	Folders []FolderSpec
}

// Vault
type TLockVault struct {
	// Path to the vault file
	Path string

	// Data
	Data TLockVaultData

	// Password to encrypt with
	password string
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

// == Vault actions ==

// == Folders functions ==

// Adds a new folder with `name`
func (vault *TLockVault) AddFolder(name string) {
	vault.Data.Folders = append(vault.Data.Folders, FolderSpec{Name: name})

	vault.write()
}

// Adds a new folder with `name`
func (vault *TLockVault) RenameFolder(old_name, new_name string) {
	folder_index := vault.find_folder(old_name)

	vault.Data.Folders[folder_index].Name = new_name

	vault.write()
}

// Returns all the tokens inside of a folder
func (vault *TLockVault) GetTokens(folder string) []TokenURI {
	folder_index := vault.find_folder(folder)

	return vault.Data.Folders[folder_index].Uris
}

// Deletes a folder with the given name
func (vault *TLockVault) DeleteFolder(name string) {
	vault.Data.Folders = remove(vault.Data.Folders, vault.find_folder(name))

	vault.write()
}

// Adds a new token to the folder
func (vault *TLockVault) AddToken(folder, uri string) {
	folder_index := vault.find_folder(folder)

	vault.Data.Folders[folder_index].Uris = append(vault.Data.Folders[folder_index].Uris, TokenURI{URI: uri, UsageCounter: 0})

	vault.write()
}

// Moves the folder up
func (vault *TLockVault) MoveFolderUp(folder string) bool {
	folder_index := vault.find_folder(folder)

    // If is folder is already at top, just return; we dont need to do anything
    if folder_index == 0 {
        return false
    }

    vault.Data.Folders[folder_index], vault.Data.Folders[folder_index - 1] = vault.Data.Folders[folder_index - 1], vault.Data.Folders[folder_index]

	vault.write()

    return true
}

// Moves the folder down
func (vault *TLockVault) MoveFolderDown(folder string) bool {
	folder_index := vault.find_folder(folder)

    // If is folder is already at top, just return; we dont need to do anything
    if folder_index == len(vault.Data.Folders) - 1 {
        return false
    }

    vault.Data.Folders[folder_index], vault.Data.Folders[folder_index + 1] = vault.Data.Folders[folder_index + 1], vault.Data.Folders[folder_index]

	vault.write()

    return true
}

// Increaments the usage counter for the given URI inside of a folder
func (vault *TLockVault) IncrementTokenUsageCounter(folder, target_uri string) {
    folder_index := vault.find_folder(folder)
    token_index := vault.find_token(folder_index, target_uri)

    vault.Data.Folders[folder_index].Uris[token_index].UsageCounter += 1

    vault.write()
}

// Edits a new token in the folder
func (vault *TLockVault) EditToken(folder, old, new string) {
	folder_index := vault.find_folder(folder)
    token_index := vault.find_token(folder_index, old)

    vault.Data.Folders[folder_index].Uris[token_index].URI = new

	vault.write()
}

func (vault *TLockVault) MoveURI(folder, token_uri, toFolder string) int {
	folder_index := vault.find_folder(folder)
    to_folder := vault.find_folder(toFolder)
    token_index := vault.find_token(folder_index, token_uri)

    toMove := vault.Data.Folders[folder_index].Uris[token_index]

    vault.Data.Folders[folder_index].Uris = remove(vault.Data.Folders[folder_index].Uris, token_index)
    vault.Data.Folders[to_folder].Uris = append(vault.Data.Folders[to_folder].Uris, toMove)

    vault.write()

    return 1
};

// Moves the token in a folder up
func (vault *TLockVault) MoveTokenUp(folder, token string) bool {
	folder_index := vault.find_folder(folder)
    token_index := vault.find_token(folder_index, token)

    // If is folder is already at top, just return; we dont need to do anything
    if token_index == 0 {
        return false
    }

    vault.Data.Folders[folder_index].Uris[token_index], vault.Data.Folders[folder_index].Uris[token_index - 1] = vault.Data.Folders[folder_index].Uris[token_index - 1], vault.Data.Folders[folder_index].Uris[token_index]

	vault.write()

    return true
}

// Moves the token in a folder down
func (vault *TLockVault) MoveTokenDown(folder, token string) bool {
	folder_index := vault.find_folder(folder)
    token_index := vault.find_token(folder_index, token)

    // If is folder is already at bottom, just return; we dont need to do anything
    if token_index == len(vault.Data.Folders[folder_index].Uris) - 1 {
        return false
    }

    vault.Data.Folders[folder_index].Uris[token_index], vault.Data.Folders[folder_index].Uris[token_index + 1] = vault.Data.Folders[folder_index].Uris[token_index + 1], vault.Data.Folders[folder_index].Uris[token_index]

	vault.write()

    return true
}

// Returns the index of the folder based on the name
func (vault TLockVault) find_folder(name string) int {
	return slices.IndexFunc(vault.Data.Folders, func(item FolderSpec) bool { return item.Name == name })
}

// Returns the index of the folder based on the name
func (vault TLockVault) find_token(folder_index int, token string) int {
	return slices.IndexFunc(vault.Data.Folders[folder_index].Uris, func(item TokenURI) bool { return item.URI == token })
}

