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
    // Data
    Data TLockVaultData

    // Path to the vault
    VaultPath string

    // Password for the vault
    password string
}

// Initializes a new instance of the vault for the new user with the given password
func Initialize(password string) TLockVault {
    // New uuid
    id := uuid.New()
    dir := path.Join(xdg.DataHome, "tlock", "root", id.String())

    // Make root dir
    if err := os.MkdirAll(dir, os.ModePerm); err != nil {
        log.Fatalf("Failed to create user's root dir: %v", err)
    }

    // Initialize vault
    vault := TLockVault {
        Data: TLockVaultData{},
        password: password,
        VaultPath: path.Join(dir, "vault.dat"),
    }

    // Write empty data
    vault.write()

    // Return
    return vault
}

// Loads the vault at the given location
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
    return &TLockVault {
        Data: data,
        VaultPath: path,
        password: password,
    }, nil
}

// [PRIVATE] Writes the current data to the file by encrypting data
func (vault TLockVault) write() {
    // Serialize
    serialized, _ := binary.Marshal(vault.Data)

    // Encrypt
    encrypted := Encrypt(vault.password, serialized);

    // Create file
    f, _ := os.Create(vault.VaultPath)

    // Write
    if _, err := f.Write(encrypted); err != nil {
        log.Fatalf("Failed to write to file: %v", err)
    }
}

// Adds a new folder with `name`
func (vault *TLockVault) AddFolder(name string) {
    vault.Data.Folders = append(vault.Data.Folders, FolderSpec{ Name: name })

    vault.write()
}

// Deletes a folder with the given name
func (vault *TLockVault) DeleteFolder(name string) {
    vault.Data.Folders = remove(vault.Data.Folders, vault.find_folder(name))

    vault.write()
}

// Adds a new URI to the given folder name
func (vault *TLockVault) AddURI(folder, uri string) {
    index := vault.find_folder(folder)

    vault.Data.Folders[index].Uris = append(vault.Data.Folders[index].Uris, uri)

    vault.write()
}

// Updates a single uri with new value
func (vault *TLockVault) UpdateURI(folder, uri, newURI string) {
    folder_index, uri_index := vault.find_folder_and_uri(folder, uri)

    vault.Data.Folders[folder_index].Uris[uri_index] = newURI

    vault.write()
}

// Deletes a URI
func (vault *TLockVault) DeleteURI(folder, uri string) {
    folder_index, uri_index := vault.find_folder_and_uri(folder, uri)

    vault.Data.Folders[folder_index].Uris = remove(vault.Data.Folders[folder_index].Uris, uri_index)

    vault.write()
}

// Moves down a uri
func (vault *TLockVault) MoveDown(folder, uri string) int {
    folder_index, uri_index := vault.find_folder_and_uri(folder, uri)

    // Uri is already at the bottom most
    if uri_index == len(vault.Data.Folders[folder_index].Uris) - 1 {
        return 0
    }

    // Swap
    vault.swap(folder_index, uri_index, uri_index + 1)

    // Write
    vault.write()

    return 1
}

// Moves down a uri
func (vault *TLockVault) MoveUp(folder, uri string) int {
    folder_index, uri_index := vault.find_folder_and_uri(folder, uri)

    // Uri is already at the top most
    if uri_index == 0 {
        return 0
    }

    // Swap
    vault.swap(folder_index, uri_index, uri_index - 1)

    // Write
    vault.write()

    return 1
}

// Moves a uri from a folder to another folder
func (vault *TLockVault) MoveURI(uri, folder, toFolder string) {
    folder_index, uri_index := vault.find_folder_and_uri(folder, uri)

    // Remove from the existing folder
    vault.Data.Folders[folder_index].Uris = remove(vault.Data.Folders[folder_index].Uris, uri_index)

    // Add to the folder
    // Add URI will handle the write
    vault.AddURI(toFolder, uri)
}

// Returns the index of the folder based on the name
func (vault TLockVault) find_folder(name string) int {
    return slices.IndexFunc(vault.Data.Folders, func(item FolderSpec) bool { return item.Name == name })
}

// Returns the index of the folder based on the name and the uri
func (vault TLockVault) find_folder_and_uri(folder, uri string) (int, int) {
    if folder_index := vault.find_folder(folder); folder_index != -1 {
        if uri_index := slices.Index(vault.Data.Folders[folder_index].Uris, uri); uri_index != -1 {
            return folder_index, uri_index
        }
    }

    return -1, -1
}

// Swaps two URI index
func (vault TLockVault) swap(folder_index, uri_index1, uri_index2 int) {
    // Classic swap
    temp := vault.Data.Folders[folder_index].Uris[uri_index1]
    vault.Data.Folders[folder_index].Uris[uri_index1] = vault.Data.Folders[folder_index].Uris[uri_index2]
    vault.Data.Folders[folder_index].Uris[uri_index2] = temp
}

