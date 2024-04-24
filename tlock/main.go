package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/muesli/termenv"

	tea "github.com/charmbracelet/bubbletea"
	tlockmodels "github.com/eklairs/tlock/tlock-models"
)

// Patches the stdout based on the context
// Patching basically changes the background color to the theme's Background color
func PatchOutput(context context.Context) (termenv.Color, io.Writer) {
	r, g, b, _ := context.GetCurrentTheme().Background.RGBA()

	// New termenv output
	output := termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.TrueColor))

	// Save background color
	original_bg := output.BackgroundColor()

	// Set background color
	output.SetBackgroundColor(output.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b)))

	// Return
	return original_bg, output
}

// Undos the patch and bring backs the original color
func ClearOutput(original_bg termenv.Color, output io.Writer) {
	term_output, ok := output.(*termenv.Output)

	if !ok {
		panic("Invalid output recieved - expected a termenv output")
	}

	// Reset
	term_output.SetBackgroundColor(original_bg)
}

// TLock go brrr
func main() {
	// Initialize context
	context := context.InitializeContext()

	original_bg, output := PatchOutput(context)

	// Start program
	program := tea.NewProgram(tlockmodels.InitializeRootModel(), tea.WithAltScreen(), tea.WithOutput(output))

	// Run
	if _, err := program.Run(); err != nil {
		log.Fatalf("[tlock] Error while running app: %v", err)
	}

	// Clear
	ClearOutput(original_bg, output)
}
