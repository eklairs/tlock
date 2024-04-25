package dashboard

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	"github.com/muesli/termenv"
)

const (
	stateFiltering = iota
	stateChoosing
)

// Themes key map
type themesKeyMap struct {
	Esc    key.Binding
	Save   key.Binding
	Up key.Binding
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
type themeItem tlockstyles.Theme

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

	fmt.Fprint(w, render_fn(m.Width() - 6, item.Name, ""))
}

var themesAsciiArt = `
▀█▀ █ █ █▀▀ █▀▄▀█ █▀▀
 █  █▀█ ██▄ █ ▀ █ ██▄`

// Themes screen
type ThemesScreen struct {
	// Context
	context context.Context

	// List
	listview list.Model

	// Help
	help help.Model

	// State
	state int

	// Filter
	filter textinput.Model
}

// Initializes a new instance of the themes screen
func InitializeThemesScreen(context context.Context) ThemesScreen {
	themeItems := make([]list.Item, len(context.Themes))

	for index, theme := range context.Themes {
		themeItems[index] = themeItem(theme)
	}

    filter := components.InitializeInputBox("Search for theme...")
    filter.Focus()

	return ThemesScreen{
		context:  context,
		listview: components.ListViewSimple(themeItems, themeListDelegate{}, 65, 18),
		help:     components.BuildHelp(),
		state:    stateChoosing,
		filter:   filter,
	}
}

// Init
func (screen ThemesScreen) Init() tea.Cmd {
	return nil
}

// Update
func (screen ThemesScreen) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
    cmds := make([]tea.Cmd, 0)

    previousFilterValue := screen.filter.Value()

    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msgType, themesKeys.Up) || key.Matches(msgType, themesKeys.Down):
        case key.Matches(msgType, themesKeys.Save):
            newTheme := screen.listview.Items()[screen.listview.Index()].(themeItem)

            screen.context.SetTheme(newTheme.Name)

            manager.PopScreen()
        default:
            screen.filter, _ = screen.filter.Update(msg)

            themeItems := make([]list.Item, 0)

            for _, theme := range screen.context.Themes {
                if strings.Index(theme.Name, screen.filter.Value()) != -1 {
                    themeItems = append(themeItems, themeItem(theme))
                }
            }

            screen.listview.SetItems(themeItems)
        }
    }

    // Update theme if the list item was changed
    previousIndex := screen.listview.Index()
	screen.listview, _ = screen.listview.Update(msg)

    if previousIndex != screen.listview.Index() || previousFilterValue != screen.filter.Value() {
        // Get new theme
        newTheme := screen.listview.Items()[screen.listview.Index()].(themeItem)

        // Change background color
        cmds = append(cmds, func() tea.Msg {
            return tea.SetBackgroundColor(termenv.RGBColor(newTheme.Background))
        })

        // Reinitialize styles
        tlockstyles.InitializeStyles(tlockstyles.Theme(newTheme))
    }

	return screen, tea.Batch(cmds...)
}

// View
func (screen ThemesScreen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		tlockstyles.Styles.Title.Render(themesAsciiArt), "",
		tlockstyles.Styles.SubText.Render("Choose a theme for tlock"), "",
        tlockstyles.Styles.Input.Render(screen.filter.View()), "",
		screen.listview.View(), "",
		screen.help.View(themesKeys),
	)
}
