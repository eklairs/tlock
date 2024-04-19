package modelmanager

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// A tea.Model-ish interface but for model manager
type Screen interface {
	Init() tea.Cmd
	Update(tea.Msg, *ModelManager) (Screen, tea.Cmd)
	View() string
}

type stubMessage struct{}

// A quick and simple utility that can effectively manage mutliple models
// NOTE: I believe this is enough for simple use cases, but for complex, I am not sure if this is the better way
type ModelManager struct {
	screens []Screen
	width   int
	height  int
}

// Initializes a new instance of model manager
func New(rootScreen Screen) ModelManager {
	return ModelManager{
		screens: []Screen{rootScreen},
	}
}

// Initializes the root model
func (manager ModelManager) Init() tea.Cmd {
	return manager.screens[0].Init()
}

// Update
func (manager ModelManager) Update(msg tea.Msg) (ModelManager, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	// React to the message
	switch msgType := msg.(type) {
	case tea.WindowSizeMsg:
		manager.width = msgType.Width
		manager.height = msgType.Height
	}

	// Update current model
	var cmd tea.Cmd
	currentScreenIndex := len(manager.screens) - 1

	manager.screens[currentScreenIndex], cmd = manager.screens[currentScreenIndex].Update(msg, &manager)

	cmds = append(cmds, cmd)

	// Return
	return manager, tea.Batch(cmds...)
}

// View()
func (manager ModelManager) View() string {
	return lipgloss.Place(
		manager.width, manager.height,
		lipgloss.Center, lipgloss.Center,
		manager.screens[len(manager.screens)-1].View(),
	)
}

// Pushes a new screen to backstack
func (manager *ModelManager) PushScreen(screen Screen) tea.Cmd {
	manager.screens = append(manager.screens, screen)

	return screen.Init()
}

// Pops the screen from the backstack
func (manager *ModelManager) PopScreen() tea.Cmd {
	if len(manager.screens) > 1 {
		manager.screens = manager.screens[:len(manager.screens)-1]
	}

	return func() tea.Msg {
		return stubMessage{}
	}
}

