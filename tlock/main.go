package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/eklairs/tlock/models"
)

// TLock go brrr
func main() {
    program := tea.NewProgram(models.InitializeRootModel(), tea.WithAltScreen());

    if _, err := program.Run(); err != nil {
        log.Fatalf("Error while running tlock: %v", err)
    }
}

