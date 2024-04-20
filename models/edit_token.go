package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/internal/modelmanager"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"github.com/pquerna/otp/totp"
	"github.com/pquerna/otp"
)

var ____ascii = `
 _____   _ _ _   
|   __|_| |_| |_ 
|   __| . | |  _|
|_____|___|_|_|  
`

type editTokenStyles struct {
    title lipgloss.Style
    titleCenter lipgloss.Style
    dimmed lipgloss.Style
    dimmedCenter lipgloss.Style
    input lipgloss.Style
}

// Root Model
type EditTokenModel struct {
    styles editTokenStyles
    vault *tlockvault.TLockVault
    folder int
    original int
    inputs []textinput.Model
    inputFocusIndex int
}

func buildTextInput(width int, value, placeholder string, styles editTokenStyles) textinput.Model {
    inputBox := textinput.New();
    inputBox.Prompt = ""
    inputBox.Width = width
    inputBox.SetValue(value)
    inputBox.Placeholder = placeholder
    inputBox.PlaceholderStyle = styles.dimmed.Copy().Background(lipgloss.Color("#1e1e2e"))

    return inputBox
}

// Initialize root model
func InitializeEditTokenModel(vault *tlockvault.TLockVault, folder, original int) EditTokenModel {
    dimmed := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
    styles := editTokenStyles{
        title: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")),
        titleCenter: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")).Width(55).Align(lipgloss.Center),
        input: lipgloss.NewStyle().Padding(1, 3).Width(55).Background(lipgloss.Color("#1e1e2e")),
        dimmed: dimmed,
        dimmedCenter: dimmed.Width(55).Copy().Align(lipgloss.Center),
    }

    token, _ := otp.NewKeyFromURL(vault.Data.Folders[folder].Uris[original])

    nameInput := buildTextInput(58, token.AccountName(), "Enter account's name...", styles)
    issuerInput := buildTextInput(58, token.Issuer(), "Enter issuer's name...", styles)
    secretInput := buildTextInput(58, token.Secret(), "Enter the secret sauce...", styles)

    nameInput.Focus()

    return EditTokenModel {
        styles: styles,
        vault: vault,
        folder: folder,
        original: original,
        inputs: []textinput.Model { nameInput, issuerInput, secretInput },
        inputFocusIndex: 0,
    }
}

// Init
func (m EditTokenModel) Init() tea.Cmd {
    return nil
}

// Update
func (m EditTokenModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch msgType.String() {
        case "tab":
            m.inputFocusIndex = (m.inputFocusIndex + 1) % 3
        case "shift+tab":
            if m.inputFocusIndex == 0 {
                m.inputFocusIndex = 2
            } else {
                m.inputFocusIndex -= 1
            }
        case "enter":
            key, _ := totp.Generate(totp.GenerateOpts{
                Issuer: m.inputs[1].Value(),
                Secret: []byte(m.inputs[2].Value()),
                AccountName: m.inputs[0].Value(),
            })

            m.vault.UpdateURI(m.folder, m.original, key.URL())
            manager.PopScreen()
        }
    }

    var cmd tea.Cmd

    for index, _ := range m.inputs {
        if index == m.inputFocusIndex {
            m.inputs[index].Focus()
        } else {
            m.inputs[index].Blur()
        }
    }

    m.inputs[m.inputFocusIndex], _ = m.inputs[m.inputFocusIndex].Update(msg)

	return m, cmd
}

// View
func (m EditTokenModel) View() string {
    return lipgloss.JoinVertical(
        lipgloss.Left,
        m.styles.titleCenter.Render(____ascii),
        m.styles.title.Render("Name"),
        m.styles.dimmed.Render("Name of the account, like Komaru"),
        m.styles.input.Render(m.inputs[0].View()), "",
        m.styles.title.Render("Issuer"),
        m.styles.dimmed.Render("Name of the issuer, like Telegram"),
        m.styles.input.Render(m.inputs[1].View()), "",
        m.styles.title.Render("Secret"),
        m.styles.dimmed.Render("The super secret provided by the issuer"),
        m.styles.input.Render(m.inputs[2].View()),
    )
}

