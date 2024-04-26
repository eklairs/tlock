package tlockmodels

import (
    tea "github.com/charmbracelet/bubbletea"
)

// Root model
type RootModel struct {

}

// Initializes a new instance of the root model
func InitializeRootModel() RootModel {
    return RootModel{}
}

// Init
func (model RootModel) Init() tea.Cmd {
    return nil
}

// Update
func (model RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    cmds := make([]tea.Cmd, 0)

    switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			cmds = append(cmds, tea.Quit)
		}
	}

    return model, tea.Batch(cmds...)
}

// View
func (model RootModel) View() string {
    return ""
}

