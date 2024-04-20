package models

import "github.com/charmbracelet/lipgloss"

// Colors
var COLOR_ACCENT = lipgloss.Color("4")
var COLOR_DIMMED = lipgloss.Color("8")
var COLOR_RED = lipgloss.Color("1")
var COLOR_YELLOW = lipgloss.Color("3")
var COLOR_BG_OVER = lipgloss.Color("#1E1E2E")

// Styles
type Styles struct {
    // Title
    title lipgloss.Style

    // Same as title, but with center alignment
    titleCenter lipgloss.Style

    // Dimmed
    dimmed lipgloss.Style

    // Same as dimmed, but with COLOR_BG_OVER as background color
    dimmedBgOver lipgloss.Style

    // Same as dimmed, but with center alignment
    dimmedCenter lipgloss.Style

    // For active list item
    active lipgloss.Style

    // For inactive list item
    inactive lipgloss.Style

    // Center
    center lipgloss.Style

    // Input
    input lipgloss.Style

    // Error
    error lipgloss.Style

    // Folder inactive
    folderInactive lipgloss.Style

    // Folder active
    folderActive lipgloss.Style
}

// Initializes the styles
func InitializeStyles(width int) Styles {
    // Base
    base := lipgloss.NewStyle().Width(width)

    // Title style
    title := base.Copy().Bold(true).Foreground(COLOR_ACCENT)

    // Dimmed style
    dimmed := base.Copy().Foreground(COLOR_DIMMED)

    // List item common
    listItem := base.Copy().Padding(1, 3)

    // Return
    return Styles{
        title: title,
        dimmed: dimmed,
        inactive: listItem,
        error: base.Copy().Foreground(COLOR_RED),
        folderInactive: base.Copy().Padding(1, 3),
        dimmedBgOver: dimmed.Copy().Background(COLOR_BG_OVER),
        input: base.Copy().Padding(1, 3).Background(COLOR_BG_OVER),
        center: base.Copy().Align(lipgloss.Center, lipgloss.Center),
        titleCenter: title.Copy().Align(lipgloss.Center, lipgloss.Center),
        dimmedCenter: dimmed.Copy().Align(lipgloss.Center, lipgloss.Center),
        active: listItem.Copy().Background(COLOR_BG_OVER).Foreground(COLOR_ACCENT),
        folderActive: base.Copy().Background(COLOR_BG_OVER).Padding(1, 2).BorderBackground(lipgloss.Color("#1E1E2E")).Border(lipgloss.ThickBorder(), false, false, false, true),
    }
}

