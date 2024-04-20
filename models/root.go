package models

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/eklairs/tlock/internal/modelmanager"
	tlockcore "github.com/eklairs/tlock/tlock-core"
)

// Root Model
type RootModel struct {
    modelmanager modelmanager.ModelManager
}

// Initialize root model
func InitializeRootModel() RootModel {
    core := tlockcore.New()

    return RootModel {
        modelmanager: modelmanager.New(InitializeSelectUserModel(core)),
    }
}

// Init
func (m RootModel) Init() tea.Cmd {
    return m.modelmanager.Init()
}

// Update
func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
	var cmds []tea.Cmd = make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
        switch msg.String() {
        case "q":
            cmds = append(cmds, tea.Quit)
        }
	}

    // Update model manager
    m.modelmanager, cmd = m.modelmanager.Update(msg)
    cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View
func (m RootModel) View() string {
    return m.modelmanager.View()
}

