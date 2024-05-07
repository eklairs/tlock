package tlockmessages

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

// Notifies requirement for a refresh of the folders list
type RefreshFoldersMsg struct{}

// Notifies requirement for a refresh of the focused folder token list
type RefreshTokensMsg struct{}

// Notifies folder changed
type FolderChanged struct {
	Folder tlockvault.Folder
}

// Requests to post folder changed message
type RequestFolderChanged struct{}

// Notification to update the tokens
// This is sent after every second after the dashboard has been loaded
type RefreshTokensValue struct{}

func DispatchRefreshTokensValueMsg() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return RefreshTokensValue{}
	})
}

// User has been deleted
type UserDeletedMsg struct{}

// User has been edited
type UserEditedMsg struct {
	NewName string
}
