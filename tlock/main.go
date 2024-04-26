package main

import (
	"github.com/rs/zerolog/log"

	tea "github.com/charmbracelet/bubbletea"
	tlockmodels "github.com/eklairs/tlock/tlock-models"
)

// TLock go brrr
func main() {
    // New bubbletea program
    program := tea.NewProgram(tlockmodels.InitializeRootModel(), tea.WithAltScreen())

    // Run
    if _, err := program.Run(); err != nil {
        log.Fatal().Err(err).Msg("[tlock] Error while running tlock program")
    }
}

