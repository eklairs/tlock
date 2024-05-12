package tlockvault

import (
	"slices"

	"github.com/eklairs/tlock/tlock-internal/utils"
)

// Adds a new folder to the vault
func (vault *Vault) AddFolder(name string) error {
	var err error

	// Validate
	if name, err = vault.validateFolderName(name); err == nil {
		// Add folder
		vault.Folders = append(vault.Folders, Folder{Name: name})

		// Write
		vault.write()
	}

	// Return
	return err
}

// Renames the folder to a new name
func (vault *Vault) RenameFolder(old, newName string) error {
	var err error

	// Validate folder name
	if newName, err = vault.validateFolderName(newName); err == nil {
		// Update
		vault.Folders[vault.findFolder(old)].Name = newName

		// Write
		vault.write()
	}

	// Return
	return err
}

// Returns all the tokens inside of a folder
func (vault *Vault) GetTokens(folder string) []Token {
	return vault.Folders[vault.findFolder(folder)].Tokens
}

// Deletes a folder by its name
func (vault *Vault) DeleteFolder(name string) {
	if index := vault.findFolder(name); index != -1 {
		// Remove
		vault.Folders = utils.Remove(vault.Folders, index)

		// Write
		vault.write()
	}
}

// Moves the folder up
func (vault *Vault) MoveFolderUp(name string) bool {
	// We will skip if the folder is already at top
	if index := vault.findFolder(name); index != 0 {
		// Swap
		vault.Folders = utils.Swap(vault.Folders, index, index-1)

		// Wrap
		vault.write()

		// True
		return true
	}

	return false
}

// Moves the folder down
func (vault *Vault) MoveFolderDown(name string) bool {
	// We will skip if the folder is already at bottom
	if index := vault.findFolder(name); index != len(vault.Folders)-1 {
		// Swap
		vault.Folders = utils.Swap(vault.Folders, index, index+1)

		// Wrap
		vault.write()

		// True
		return true
	}

	return false
}

// Checks if the folder with the name exists
func (vault Vault) folderExists(name string) bool {
	return vault.findFolder(name) != -1
}

// Returns the index of the folder with the given name
func (vault Vault) findFolder(name string) int {
	return slices.IndexFunc(vault.Folders, func(folder Folder) bool { return folder.Name == name })
}
