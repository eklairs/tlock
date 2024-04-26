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

// Notification to update the tokens
// This is sent after every second after the dashboard has been loaded
type RefreshTokensValue struct{}

func DispatchRefreshTokensValueMsg() tea.Msg {
	currentTime := time.Now()
	nextSecond := currentTime.Truncate(time.Second).Add(time.Second)
	duration := nextSecond.Sub(currentTime)

	time.Sleep(duration)

	return RefreshTokensValue{}
}
