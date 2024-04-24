package components

import "github.com/charmbracelet/bubbles/list"

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
