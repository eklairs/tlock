package main

import (
	"github.com/muesli/termenv"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/eklairs/tlock/tlock-internal/context"
	tlockmodels "github.com/eklairs/tlock/tlock/models"
	tlockstyles "github.com/eklairs/tlock/tlock/styles"
)

// TLock go brrr
func main() {
	// Initialize context
	context := context.InitializeContext()
	background := termenv.RGBColor(context.GetCurrentTheme().Background)

	// Initialize styles
	tlockstyles.InitializeStyles(context.GetCurrentTheme())

	// New bubbletea program
	program := tea.NewProgram(tlockmodels.InitializeRootModel(&context), tea.WithAltScreen(), tea.WithBackgroundColor(background), tea.WithMouseCellMotion())

	// Run
	if _, err := program.Run(); err != nil {

	}
}
