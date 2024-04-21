package models

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Root Model
type RootModel struct {

}

// Initialize root model
func InitializeRootModel() RootModel {
	return RootModel{

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

	return m, tea.Batch(cmds...)
}

// View
func (m RootModel) View() string {
    return ""
}

