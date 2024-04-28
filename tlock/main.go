package main

import (
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/muesli/termenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/utils"
	tlockmodels "github.com/eklairs/tlock/tlock-models"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

var LOG_FILE = path.Join(xdg.DataHome, "tlock", "logs", "log")

// Returns the log writer for zerolog, which points to the log file
func setupZeroLog() *os.File {
	// Use unix timestamp
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Create the log file
	file, err := utils.EnsureExists(LOG_FILE)

	// If error, disable logs
	if err != nil {
		return nil
	}

	// Set logger
	log.Logger = zerolog.New(file)

	// Return file
	return file
}

// TLock go brrr
func main() {
	// Setup zerolog
	if file := setupZeroLog(); file != nil {
		defer file.Close()
	}

	// Initialize context
	context := context.InitializeContext()
	background := termenv.RGBColor(context.GetCurrentTheme().Background)

	// Initialize styles
	tlockstyles.InitializeStyles(context.GetCurrentTheme())

	// New bubbletea program
	program := tea.NewProgram(tlockmodels.InitializeRootModel(&context), tea.WithAltScreen(), tea.WithBackgroundColor(background))

	// Run
	if _, err := program.Run(); err != nil {
		log.Fatal().Err(err).Msg("[tlock] Error while running tlock program")
	}
}
