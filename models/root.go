package models

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/eklairs/tlock/internal/modelmanager"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
)

// Root Model
type RootModel struct {
    modelmanager modelmanager.ModelManager
}

// Initialize root model
func InitializeRootModel() RootModel {
    temp_vault, _ := tlockvault.Load("/home/kyeboard/.local/share/tlock/root/6153db12-9995-4c50-abb4-584be0216550/vault.dat", "")

    return RootModel {
        modelmanager: modelmanager.New(InitializeNewFolderModel(*temp_vault)),
    }
}

// Init
func (m RootModel) Init() tea.Cmd {
    return m.modelmanager.Init()
}

// Update
func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
        switch msg.String() {
        case "q":
            cmds = append(cmds, tea.Quit)
        }
	}

    m.modelmanager, _ = m.modelmanager.Update(msg)

	return m, tea.Batch(cmds...)
}

// View
func (m RootModel) View() string {
    return m.modelmanager.View()
}

