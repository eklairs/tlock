package main

import (
	"fmt"
	"log"

	"github.com/eklairs/tlock/tlock-internal/context"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	"github.com/muesli/termenv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tlockmodels "github.com/eklairs/tlock/tlock-models"
)

// Changes lipgloss.color to termenv.Color
func ToTermenvColor(color lipgloss.Color) termenv.Color {
	r, g, b, _ := color.RGBA()

	// Set background color
	return termenv.RGBColor(fmt.Sprintf("#%02x%02x%02x", r, g, b))
}

// TLock go brrr
func main() {
	// Initialize context
	context := context.InitializeContext()

	// Initialize styles
	tlockstyles.InitializeStyles(context.GetCurrentTheme())

	// Start program
	program := tea.NewProgram(tlockmodels.InitializeRootModel(context), tea.WithAltScreen(), tea.WithBackgroundColor(ToTermenvColor(context.GetCurrentTheme().Background)))

	// Run
	if _, err := program.Run(); err != nil {
		log.Fatalf("[tlock] Error while running app: %v", err)
	}
}
