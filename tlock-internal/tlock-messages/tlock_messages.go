package tlockmessages

import tlockvault "github.com/eklairs/tlock/tlock-vault"

// Notifies requirement for a refresh of the folders list
type RefreshFoldersMsg struct{}

// Notifies requirement for a refresh of the focused folder token list
type RefreshTokensMsg struct{}

// Notifies folder changed
type FolderChanged struct {
	Folder tlockvault.Folder
}
