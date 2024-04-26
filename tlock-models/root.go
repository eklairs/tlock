package tlockmodels

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eklairs/tlock/tlock-internal/context"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	tlockmessages "github.com/eklairs/tlock/tlock-internal/tlock-messages"
	"github.com/eklairs/tlock/tlock-models/auth"
)

// Root model
type RootModel struct {
	manager modelmanager.ModelManager
}

// Initializes a new instance of the root model
func InitializeRootModel(context context.Context) RootModel {
	return RootModel{
		manager: modelmanager.New(auth.InitializeSelectUserScreen(context)),
	}
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
	// We dispatch back the message from root model because its the only model that recieves all the models everytime.
	// If a new screen is pushed to modelmanager, the dashboard will not recieve the message and thus will break the update
	case tlockmessages.RefreshTokensValue:
		cmds = append(cmds, tlockmessages.DispatchRefreshTokensValueMsg)
	}

	// Update model manager
	cmds = append(cmds, model.manager.Update(msg))

	// Return
	return model, tea.Batch(cmds...)
}

// View
func (model RootModel) View() string {
	return model.manager.View()
}
