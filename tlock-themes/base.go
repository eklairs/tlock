package tlockthemes

import "github.com/charmbracelet/lipgloss"

// Theme object
type Theme struct {
	// Main background
	WindowBg lipgloss.Color

	// Main foregorund
	WindowFg lipgloss.Color

	// Background over
	WindowBgOver lipgloss.Color

	// Foreground over
	WindowFgOver lipgloss.Color

	// Dimmed color
	Dimmed lipgloss.Color

	// Accent color
	Accent lipgloss.Color

	// Error
	Error lipgloss.Color
}

