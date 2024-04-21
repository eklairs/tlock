package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-models/login"
)

// Root Model
type RootModel struct {
    modelmanager modelmanager.ModelManager
}

// Initialize root model
func InitializeRootModel() RootModel {
    context := context.InitializeContext()

    // Screen to initialize with
	var screen modelmanager.Screen

	if len(context.Core.Users) == 0 {
		screen = login.InitializeCreateUserModel(context)
	} else {
		screen = login.InitializeSelectUserModel(context)
	}

    return RootModel{
        modelmanager: modelmanager.New(screen),
    }
}

// Init
func (m RootModel) Init() tea.Cmd {
	return nil
}

// Update
func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			cmds = append(cmds, tea.Quit)
		}
	}

    // Update model manager
	cmds = append(cmds, m.modelmanager.Update(msg))

	return m, tea.Batch(cmds...)
}

// View
func (m RootModel) View() string {
    return m.modelmanager.View()
}

