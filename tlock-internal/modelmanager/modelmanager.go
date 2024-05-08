package modelmanager

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"golang.org/x/term"
)

// Message that represents the model has come back to focus
type ScreenRefocusedMsg struct{}

// A tea.Model-ish interface but for model manager
type Screen interface {
	Init() tea.Cmd
	Update(tea.Msg, *ModelManager) (Screen, tea.Cmd)
	View() string
}

// Operations
const (
	OperationNone = iota
	OperationPush
	OperationPop
	OperationReplace
)

// Type of operation
type Operation struct {
	// Action
	Action int

	// Related screen
	Screen *Screen
}

// None operation
func NoneOperation() Operation {
	return Operation{
		Action: OperationNone,
		Screen: nil,
	}
}

// Model manager
type ModelManager struct {
	// Backstack
	stack []Screen

	// Any pending operation
	operation Operation
}

// Initializes a new instance of the model manager
func New(rootScreen Screen) ModelManager {
	return ModelManager{
		stack:     []Screen{rootScreen},
		operation: NoneOperation(),
	}
}

// Adds a new screen on the stack
func (manager *ModelManager) PushScreen(screen Screen) tea.Cmd {
	manager.operation = Operation{
		Action: OperationPush,
		Screen: &screen,
	}

	return screen.Init()
}

// Replaces the current screen with a new screen on the stack
func (manager *ModelManager) ReplaceScreen(screen Screen) tea.Cmd {
	manager.operation = Operation{
		Action: OperationReplace,
		Screen: &screen,
	}

	return screen.Init()
}

// Pops the top screen from the stack
func (manager *ModelManager) PopScreen() {
	if len(manager.stack) > 1 {
		manager.operation = Operation{
			Action: OperationPop,
			Screen: nil,
		}
	}
}

// Calls the update method on the current screen
func (manager *ModelManager) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	// Current screen is the screen at the top
	screen_index := len(manager.stack) - 1

	// Update
	switch msg.(type) {
	case components.StatusBarMsg:
		// Send it to all the screens
		for i := 0; i < len(manager.stack); i++ {
			manager.stack[i], cmd = manager.stack[i].Update(msg, manager)
		}
	default:
		manager.stack[screen_index], cmd = manager.stack[screen_index].Update(msg, manager)
	}

	// Resolve operation
	manager.ResolveOperation()

	return cmd
}

// Resolves the operation
func (manager *ModelManager) ResolveOperation() {
	// Current screen is the screen at the top
	screen_index := len(manager.stack) - 1

	// Resolve any pending operation
	switch manager.operation.Action {
	case OperationPush:
		manager.stack = append(manager.stack, *manager.operation.Screen)
	case OperationPop:
		manager.stack = manager.stack[:screen_index]

		// Send new message
		manager.Update(ScreenRefocusedMsg{})

	case OperationReplace:
		manager.stack[screen_index] = *manager.operation.Screen
	}

	// Reset operation
	manager.operation = NoneOperation()
}

// Calls the View() function on the current screen with center aligned to the screen
func (manager ModelManager) View() string {
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		manager.stack[len(manager.stack)-1].View(),
	)
}
