package main

import (
	"log"

	tlockinternal "github.com/eklairs/tlock/tlock-internal"
	"github.com/eklairs/tlock/tlock-internal/context"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"

	tea "github.com/charmbracelet/bubbletea"
	tlockmodels "github.com/eklairs/tlock/tlock-models"
)

// TLock go brrr
func main() {
	// Initialize context
	context := context.InitializeContext()

	// Initialize styles
	tlockstyles.InitializeStyles(context.GetCurrentTheme())

	// Start program
	program := tea.NewProgram(tlockmodels.InitializeRootModel(context), tea.WithAltScreen(), tea.WithBackgroundColor(tlockinternal.ToTermenvColor(context.GetCurrentTheme().Background)))

	// Run
	if _, err := program.Run(); err != nil {
		log.Fatalf("[tlock] Error while running app: %v", err)
	}
}
