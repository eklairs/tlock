package tlockinternal

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// Changes lipgloss.color to termenv.Color
func ToTermenvColor(color lipgloss.Color) termenv.Color {
	r, g, b, _ := color.RGBA()

	// Set background color
	return termenv.RGBColor(fmt.Sprintf("#%02x%02x%02x", r, g, b))
}

