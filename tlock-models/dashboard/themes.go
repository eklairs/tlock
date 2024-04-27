package dashboard

import (
	"fmt"
	"io"
	"slices"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tlockinternal "github.com/eklairs/tlock/tlock-internal"
	"github.com/eklairs/tlock/tlock-internal/components"
	tlockcontext "github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	"github.com/muesli/termenv"
)

// Themes key map
type themesKeyMap struct {
	Esc  key.Binding
	Save key.Binding
	Up   key.Binding
	Down key.Binding
}

// ShortHelp()
func (k themesKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Esc, k.Save}
}

// FullHelp()
func (k themesKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

// Keys
var themesKeys = themesKeyMap{
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
	Save: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "save"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
}

// Theme item
type themeItem tlockcontext.Theme

// FilterValue()
func (item themeItem) FilterValue() string {
	return item.Name
}

// Delegate
type themeListDelegate struct{}

// Height
func (d themeListDelegate) Height() int {
	return 3
}

// Spacing
func (d themeListDelegate) Spacing() int {
	return 0
}

// Update
func (d themeListDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

// Render
func (d themeListDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item := listItem.(themeItem)

	// Decide renderer function
	render_fn := components.ListItemInactive

	if index == m.Index() {
		render_fn = components.ListItemActive
	}

	fmt.Fprint(w, render_fn(m.Width()-6, item.Name, ""))
}

var themesAsciiArt = `
▀█▀ █ █ █▀▀ █▀▄▀█ █▀▀
 █  █▀█ ██▄ █ ▀ █ ██▄`

// Themes screen
type ThemesScreen struct {
	// Context
	context tlockcontext.Context

	// List
	listview list.Model
}

// Initializes a new instance of the themes screen
func InitializeThemesScreen(context tlockcontext.Context) ThemesScreen {
	// Theme items
	themeItems := tlockinternal.Map(context.Themes, func(theme tlockcontext.Theme) list.Item { return themeItem(theme) })

	// Initialize theme list
	listview := components.ListViewSimple(themeItems, themeListDelegate{}, 65, min(18, len(context.Themes)*3))

	// Set the focus to the currently applied theme
	for i := 0; i < slices.IndexFunc(context.Themes, func(t tlockcontext.Theme) bool { return t.Name == context.Config.CurrentTheme }); i++ {
		listview.CursorDown()
	}

	return ThemesScreen{
		context:  context,
		listview: listview,
	}
}

// Init
func (screen ThemesScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen ThemesScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msgType, themesKeys.Esc):
			// Get original theme
			originalTheme := screen.context.GetCurrentTheme()

			// Reset back the theme
			cmds = append(cmds, func() tea.Msg {
				return tea.SetBackgroundColor(termenv.RGBColor(originalTheme.Background))
			})

			// Reinitialize styles
			tlockstyles.InitializeStyles(tlockcontext.Theme(originalTheme))

			// Pop screen
			manager.PopScreen()

		case key.Matches(msgType, themesKeys.Save):
			// Get the focused theem
			newTheme := screen.listview.SelectedItem().(themeItem)

			// Set the theme
			screen.context.SetTheme(newTheme.Name)

			// Pop
			manager.PopScreen()
		}
	}

	// Save the previous position before updating
	previousIndex := screen.listview.Index()

	// Update listview
	screen.listview, _ = screen.listview.Update(msg)

	// Check for theme updates
	if previousIndex != screen.listview.Index() {
		// Get new theme
		newTheme := screen.listview.SelectedItem().(themeItem)

		// Change background color
		cmds = append(cmds, func() tea.Msg {
			return tea.SetBackgroundColor(termenv.RGBColor(newTheme.Background))
		})

		// Reinitialize styles
		tlockstyles.InitializeStyles(tlockcontext.Theme(newTheme))
	}

	// Return
	return screen, tea.Batch(cmds...)
}

// View
func (screen ThemesScreen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(themesAsciiArt), "",
		tlockstyles.Styles.SubText.Render("Choose a theme for tlock"), "",
		screen.listview.View(), "",
		tlockstyles.Help.View(themesKeys),
	)
}
