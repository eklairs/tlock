package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
)

// List item active
func ListItemActive(width int, title, suffix string) string {
	space_width := width - lipgloss.Width(title) - lipgloss.Width(suffix)

	ui := lipgloss.JoinHorizontal(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(title),
		tlockstyles.Styles.BackgroundOver.Render(strings.Repeat(" ", space_width)),
		tlockstyles.Styles.BackgroundOver.Render(tlockstyles.Styles.Title.Render(suffix)),
	)

	return tlockstyles.Styles.ListItemActive.Render(ui)
}

// List item active
func ListItemInactive(width int, title, suffix string) string {
	space_width := width - lipgloss.Width(title) - lipgloss.Width(suffix)

	ui := lipgloss.JoinHorizontal(
		lipgloss.Center,
		tlockstyles.Styles.SubText.Render(title),
		strings.Repeat(" ", space_width),
		tlockstyles.Styles.SubText.Render(suffix),
	)

	return tlockstyles.Styles.ListItemInactive.Render(ui)
}

// Token list item renderer implementation
func tokenItemImpl(width int, account, separator, issuer, code string, spacerStyle lipgloss.Style, uiStyle lipgloss.Style) string {
	space_width := width - lipgloss.Width(account) - lipgloss.Width(separator) - lipgloss.Width(issuer) - lipgloss.Width(code)

	var ui string

	// If the space width is > 0, we have some space for the padding!
	if space_width >= 0 {
		ui = lipgloss.JoinHorizontal(
			lipgloss.Left,
			account, separator, issuer,
			spacerStyle.Render(strings.Repeat(" ", space_width)),
			code,
		)
	} else if newSpaceWidth := space_width + lipgloss.Width(issuer) + lipgloss.Width(separator); newSpaceWidth >= 0 {
		// If the width is not enough; lets drop the issuer name
		ui = lipgloss.JoinHorizontal(
			lipgloss.Left,
			account,
			spacerStyle.Render(strings.Repeat(" ", newSpaceWidth)),
			code,
		)
	} else {
		ui = lipgloss.JoinHorizontal(lipgloss.Left, code)
	}

	// If that doesnt help, then just show the code
	return uiStyle.Render(ui)
}

// List item active
func TokenItemActive(width int, account, issuer, code string) string {
	return tokenItemImpl(
		width,
		tlockstyles.Styles.Title.Render(account),
		tlockstyles.Styles.BackgroundOver.Render(" • "),
		tlockstyles.Styles.BackgroundOver.Render(issuer),
		tlockstyles.Styles.BackgroundOver.Render(tlockstyles.Styles.Title.Render(code)),
		tlockstyles.Styles.BackgroundOver,
		tlockstyles.Styles.ListItemActive,
	)
}

// List item active
func TokenItemInactive(width int, account, issuer, code string) string {
	return tokenItemImpl(
		width,
		tlockstyles.Styles.SubText.Render(account),
		tlockstyles.Styles.SubText.Render(" • "),
		tlockstyles.Styles.SubText.Render(issuer),
		tlockstyles.Styles.SubText.Render(code),
		tlockstyles.Styles.SubText,
		tlockstyles.Styles.ListItemInactive,
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
		items = append(items, tlockstyles.Styles.Error.Copy().Width(65).Render(fmt.Sprintf("× %s", *error)), "")
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

// Active folder list item
func ActiveFolderListItem(width int, name string, tokensCount int) string {
	items := []string{
		tlockstyles.Styles.Title.Render(name),
		tlockstyles.Styles.SubText.Render(fmt.Sprintf("%d tokens", tokensCount)),
	}

	return tlockstyles.Styles.FolderItemActive.Copy().Width(width).Render(strings.Join(items, "\n"))
}

// Inactive folder list item
func InactiveFolderListItem(width int, name string, tokensCount int) string {
	items := []string{
		tlockstyles.Styles.SubText.Render(name),
		tlockstyles.Styles.SubText.Render(fmt.Sprintf("%d tokens", tokensCount)),
	}

	return tlockstyles.Styles.FolderItemInactive.Copy().Width(width).Render(strings.Join(items, "\n"))
}
