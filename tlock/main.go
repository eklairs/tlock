package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	tlockmodels "github.com/eklairs/tlock/tlock-models"
)

// TLock go brrr
func main() {
	program := tea.NewProgram(tlockmodels.InitializeRootModel(), tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		log.Fatalf("[tlock] Error while running app: %v", err)
	}
}
