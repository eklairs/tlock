package main

import (
	"log"
	"os"
	"path"

	"github.com/adrg/xdg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eklairs/tlock/models"
)

// Path to the log file
var LOG_DIR = path.Join(xdg.DataHome, "tlock", "logs")

// Sets up logging to file
func setup_logging() *os.File {
    // Make log dir
    os.MkdirAll(LOG_DIR, os.ModePerm)

    // Log to file
    file, err := tea.LogToFile(path.Join(LOG_DIR, "tlock.log"), "tlock");

    if err != nil {
        log.Fatalf("Failed to create log file at %s", LOG_DIR)
    }

    // Return
    return file
}

// TLock go brrrrr
func main() {
    // Setup logging to file
    defer setup_logging().Close()

    // Log
    log.Printf("[tlock] Starting tlock!")

    // Initialize tea program
    program := tea.NewProgram(models.InitializeRootModel(), tea.WithAltScreen());

    // Start tea program
    if _, err := program.Run(); err != nil {
        log.Fatalf("Error while running tlock: %v", err)
    }
}

