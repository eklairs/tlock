package tlockvault

import (
	"errors"
	"os"
	"path"
	"slices"

	"github.com/adrg/xdg"
	"github.com/google/uuid"
	"github.com/kelindar/binary"
	"github.com/pquerna/otp"
	"github.com/rs/zerolog/log"

	tlockinternal "github.com/eklairs/tlock/tlock-internal"
)

// Token types
const (
	TokenTypeTOTP = iota
	TokenTypeHOTP
)

// Dirs
var ROOT_DIR = path.Join(xdg.DataHome, "tlock", "root")

// Token Type
type TokenType int

// Token
type Token struct {
	// ID of the token
	ID string

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
	// ID
	ID string

	// Name of the folder
	Name string

	// Tokens
	Tokens []Token
}

// Vault securely stores all the tokens inside of the file for tlock
type Vault struct {
	// All the folders and their data
	Folders []Folder

	// Path to the file
	Path string

	// Password
	password string
}

// Initializes a new instance of the vault at the given path
func Initialize(password string) Vault {
	// Create a new folder for the new user
	id := uuid.New()
	dir := path.Join(ROOT_DIR, id.String())

	// Log if there was error while creating
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal().Err(err).Str("path", dir).Msg("Failed to create directory for user")
	}

	// Initialize vault
	vault := Vault{
		password: password,
		Path:     path.Join(dir, "vault.bin"),
	}

	// Write empty data
	vault.write()

	// Return
	return vault
}

// Writes the current data to the vault
func (vault Vault) write() {
	// Serialize
	serialized, _ := binary.Marshal(vault.Folders)

	// Encrypt
	encrypted := Encrypt(vault.password, serialized)

	// Create parent dir
	file, err := tlockinternal.EnsureExists(vault.Path)

	// Check for errors
	if err != nil {
		log.Fatal().Err(err).Str("path", vault.Path).Msg("[tlockvault] Failed to write encrypted data to file")
	}

	// Write
	file.Write(encrypted)
}

// Loads a new vault instance
// Loads a vault instance from the given path
func Load(path, password string) (*Vault, error) {
	// Read encrypted bytes
	raw, err := os.ReadFile(path)

	// No errors, pl0x
	if err != nil {
		return nil, errors.New("The vault does not exist, did you delete it?")
	}

	// Empty data
	var data []Folder

	// Decrypt
	decrypted, err := Decrypt(password, raw)

	if err != nil {
		return nil, errors.New("Error while decrypting, well that's weird")
	}

	// Unmarshal binary serialized data
	if err := binary.Unmarshal(decrypted, &data); err != nil {
		return nil, errors.New("Invalid password, please try again")
	}

	// Create vault instance and return
	return &Vault{
		password: password,
		Folders:  data,
		Path:     path,
	}, nil
}

// Adds a new folder to the vault
func (vault *Vault) AddFolder(name string) {
	// Initialize new folder
	folder := Folder{
		ID:   uuid.NewString(),
		Name: name,
	}

	// Add folder
	vault.Folders = append(vault.Folders, folder)

	// Write
	vault.write()
}

// Renames the folder to a new name
func (vault *Vault) RenameFolder(old_id, new_name string) {
	// Update
	vault.Folders[vault.find_folder(old_id)].Name = new_name

	// Write
	vault.write()
}

// Returns all the tokens inside of a folder
func (vault *Vault) GetTokens(id string) []Token {
	return vault.Folders[vault.find_folder(id)].Tokens
}

// Deletes a folder by its id
func (vault *Vault) DeleteFolder(id string) {
	// Remove folder
	vault.Folders = tlockinternal.Remove(vault.Folders, vault.find_folder(id))

	// Write
	vault.write()
}

// Moves the folder up
func (vault *Vault) MoveFolderUp(folderId string) bool {
	// Find folder index
	folder_index := vault.find_folder(folderId)

	// If is folder is already at top, just return; we dont need to do anything
	if folder_index == 0 {
		return false
	}

	// Swap
	vault.Folders = tlockinternal.Swap(vault.Folders, folder_index, folder_index-1)

	// Wrap
	vault.write()

	// Return
	return true
}

// Moves the folder down
func (vault *Vault) MoveFolderDown(folderId string) bool {
	// Find folder index
	folder_index := vault.find_folder(folderId)

	// If is folder is already at top, just return; we dont need to do anything
	if folder_index == len(vault.Folders)-1 {
		return false
	}

	// Swap
	vault.Folders = tlockinternal.Swap(vault.Folders, folder_index, folder_index+1)

	// Wrap
	vault.write()

	// Return
	return true
}

// Find a folder index by its uuid
func (vault *Vault) find_folder(id string) int {
	return slices.IndexFunc(vault.Folders, func(folder Folder) bool { return folder.ID == id })
}
