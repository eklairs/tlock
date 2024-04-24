package tlockmodels

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/eklairs/tlock/tlock-models/auth"
)

// Root Model
type RootModel struct {
	modelmanager modelmanager.ModelManager
}

// Initialize root model
func InitializeRootModel(context context.Context) RootModel {
	var screen modelmanager.Screen

	if len(context.Core.Users) == 0 {
		screen = auth.InitializeCreateUserScreen(context)
	} else {
		screen = auth.InitializeSelectUserScreen(context)
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

	cmds = append(cmds, m.modelmanager.Update(msg))

	return m, tea.Batch(cmds...)
}

// View
func (m RootModel) View() string {
	return m.modelmanager.View()
}
