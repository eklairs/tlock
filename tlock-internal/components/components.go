package components

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

// === List item active ==
func listItemImpl(width int, title, suffix string, spacerStyle, style lipgloss.Style) string {
	// Required space
	space_width := width - lipgloss.Width(title) - lipgloss.Width(suffix)

	// Join
	ui := lipgloss.JoinHorizontal(lipgloss.Center, title, spacerStyle.Render(strings.Repeat(" ", space_width)), suffix)

	// Return
	return style.Render(ui)
}

// List item active
func ListItemActive(width int, title, suffix string) string {
	return listItemImpl(
		width,
		tlockstyles.Styles.Title.Render(title),
		tlockstyles.Styles.BackgroundOver.Render(tlockstyles.Styles.Title.Render(suffix)),
		tlockstyles.Styles.BackgroundOver, tlockstyles.Styles.ListItemActive,
	)
}

// List item active
func ListItemInactive(width int, title, suffix string) string {
	return listItemImpl(
		width,
		tlockstyles.Styles.SubText.Render(title),
		tlockstyles.Styles.SubText.Render(suffix),
		tlockstyles.Styles.SubText, tlockstyles.Styles.ListItemInactive,
	)
}

// === List item active ends ==

// === Token item ==

// Token list item renderer implementation
func tokenItemImpl(width int, icon, account, separator, issuer, code string, spacerStyle lipgloss.Style, uiStyle lipgloss.Style, showIcon bool) string {
	space_width := width - lipgloss.Width(account) - lipgloss.Width(separator) - lipgloss.Width(issuer) - lipgloss.Width(code) - 3

	if showIcon {
		icon = spacerStyle.Render(fmt.Sprintf("%s%s", icon, spacerStyle.Render("  ")))
	} else {
		icon = ""
	}

	// Icon renderable
	var ui string

	// If the space width is > 0, we have some space for the padding!
	if space_width >= 0 {
		ui = lipgloss.JoinHorizontal(
			lipgloss.Left,
			icon, account, separator, issuer,
			spacerStyle.Render(strings.Repeat(" ", space_width)),
			code,
		)
	} else if newSpaceWidth := space_width + lipgloss.Width(issuer) + lipgloss.Width(separator); newSpaceWidth >= 0 {
		// If the width is not enough; lets drop the issuer name
		ui = lipgloss.JoinHorizontal(
			lipgloss.Left,
			icon, account,
			spacerStyle.Render(strings.Repeat(" ", newSpaceWidth)),
			code,
		)
	} else {
		ui = lipgloss.JoinHorizontal(lipgloss.Left, icon, code)
	}

	// If that doesnt help, then just show the code
	return uiStyle.Render(ui)
}

// List item active
func TokenItemActive(width int, icon, account, issuer, code string, period int, timeLeft *int, showIcon bool) string {
	style := tlockstyles.Styles.ListItemActive

	if timeLeft != nil {
		// We display time bar at the bottom!
		style = style.Copy().UnsetPaddingBottom()
	}

	ui := tokenItemImpl(
		width, icon,
		tlockstyles.Styles.BackgroundOver.Render(tlockstyles.Styles.Title.Render(account)),
		tlockstyles.Styles.BackgroundOver.Render(" • "),
		tlockstyles.Styles.BackgroundOver.Render(issuer),
		tlockstyles.Styles.BackgroundOver.Render(tlockstyles.Styles.Title.Render(code)),
		tlockstyles.Styles.BackgroundOver, style, showIcon,
	)

	// Get the number of blocks required to render
	if timeLeft != nil {
		// Get the number of blocks to render
		blocksPerSec := int(math.Floor(float64(width) / float64(period)))

		// Render!
		renderable := tlockstyles.Styles.Title.Render(strings.Repeat("▁", blocksPerSec**timeLeft))

		// Render!
		ui = lipgloss.JoinVertical(lipgloss.Left, ui, lipgloss.NewStyle().Inherit(style).UnsetPaddingTop().Width(width+6).Render(renderable))
	}

	return ui
}

// List item active
func TokenItemInactive(width int, icon, account, issuer, code string, period int, timeLeft *int, showIcon bool) string {
	return tokenItemImpl(
		width, icon,
		tlockstyles.Styles.SubText.Render(account),
		tlockstyles.Styles.SubText.Render(" • "),
		tlockstyles.Styles.SubText.Render(issuer),
		tlockstyles.Styles.SubText.Render(code),
		tlockstyles.Styles.SubText, tlockstyles.Styles.ListItemInactive, showIcon,
	)
}

// Builds a listview model devoid of every feature
func ListViewSimple(items []list.Item, delegate list.ItemDelegate, width, height int) list.Model {
	listview := list.New(items, delegate, width, height)

	listview.SetShowHelp(false)
	listview.SetShowTitle(false)
	listview.SetShowFilter(false)
	listview.SetShowStatusBar(false)
	listview.SetShowPagination(false)
	listview.DisableQuitKeybindings()

	return listview
}

func InputGroup(title, description string, error *string, input textinput.Model) string {
	// Total width relative to the input's width
	width := input.Width + 7

	items := []string{
		tlockstyles.Styles.Title.Copy().Width(width).Render(title),
		tlockstyles.Styles.SubText.Copy().Width(width).Render(description),
		tlockstyles.Styles.Input.Copy().Width(width).Render(input.View()), "",
	}

	// Append error if any
	if error != nil {
		items = append(items, tlockstyles.Styles.Error.Copy().Width(width).Render(fmt.Sprintf("× %s", *error)), "")
	}

	return lipgloss.JoinVertical(lipgloss.Center, items...)
}

// Creates a new input box
func InitializeInputBox(placeholder string) textinput.Model {
	return InitializeInputBoxCustomWidth(placeholder, 58)
}

func InitializeInputBoxCustomWidth(placeholder string, width int) textinput.Model {
	input := textinput.New()
	input.Prompt = ""
	input.Width = width
	input.Placeholder = placeholder
	input.PlaceholderStyle = tlockstyles.Styles.Placeholder

	return input
}

// == Folder list item
func folderListItem(width int, name string, nameStyle lipgloss.Style, tokensCount int, tokensCountStyle, style lipgloss.Style) string {
	bottom := fmt.Sprintf("%d tokens", tokensCount)

	// -6 is for including the padding
	if len(name) > width-6 {
		// Ellipsis
		name = name[:width-7] + "…"
	}

	// Wrap the tokens as well
	if len(bottom) > width-5 {
		bottom = bottom[:width-6] + "…"
	}

	// Render
	items := []string{
		nameStyle.Render(name),
		tokensCountStyle.Render(bottom),
	}

	return style.Copy().Width(width).Render(strings.Join(items, "\n"))
}

// Active folder list item
func ActiveFolderListItem(width int, name string, tokensCount int) string {
	return folderListItem(
		width, name, tlockstyles.Styles.Title,
		tokensCount, tlockstyles.Styles.SubText,
		tlockstyles.Styles.FolderItemActive,
	)
}

// Inactive folder list item
func InactiveFolderListItem(width int, name string, tokensCount int) string {
	return folderListItem(
		width, name, tlockstyles.Styles.SubText,
		tokensCount, tlockstyles.Styles.SubText,
		tlockstyles.Styles.FolderItemInactive,
	)
}
