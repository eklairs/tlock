package models

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/internal/modelmanager"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/term"
)

type TLockVendor struct {
    Icons map[string]struct {
        Unicode string
        Hex string
    }
}

type dashboardStyles struct {
    title lipgloss.Style
    titleCenter lipgloss.Style
    dimmed lipgloss.Style
    dimmedCenter lipgloss.Style
    input lipgloss.Style
}

// Root Model
type DashboardModel struct {
    styles dashboardStyles
    vault tlockvault.TLockVault
    current_index int
    token_current_index int
    vendor TLockVendor
}

var DIGIT_LIST = []string {
    `
    ┏━┓
    ┃ ┃
    ┗━┛`,
    `
     ┓
     ┃
    ╺┻╸`,
    `
    ╺━┓
    ┏━┛
    ┗━╸`,
    `
    ╺━┓
     ━┫
    ╺━┛`,
    `
    ╻ ╻
    ┗━┫
      ╹`,
    `
    ┏━╸
    ┗━┓
    ╺━┛`,
    `
    ┏━╸
    ┣━┓
    ┗━┛`,
    `
    ╺━┓
      ┃
      ╹`,
    `
    ┏━┓
    ┣━┫
    ┗━┛`,
    `
    ┏━┓
    ┗━┫
    ╺━┛`,
}

// Initialize root model
func InitializeDashboardModel(vault tlockvault.TLockVault) DashboardModel {
    dimmed := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
    raw, _ := os.ReadFile("tlock-vendor/icons.json")

    vendor := TLockVendor{}
    json.Unmarshal(raw, &vendor)

    return DashboardModel {
        styles: dashboardStyles{
            title: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")),
            titleCenter: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")).Width(30).Align(lipgloss.Center),
            input: lipgloss.NewStyle().Padding(1, 3).Width(30).Background(lipgloss.Color("#1e1e2e")),
            dimmed: dimmed,
            dimmedCenter: dimmed.Width(30).Copy().Align(lipgloss.Center),
        },
        vault: vault,
        current_index: 0,
        vendor: vendor,
    }
}

// Init
func (m DashboardModel) Init() tea.Cmd {
    return nil
}

// Update
func (m DashboardModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
    switch msgType := msg.(type) {
    case tea.KeyMsg:
        switch msgType.String() {
        case "J":
            m.current_index = (m.current_index + 1) % len(m.vault.Data.Folders)
            m.token_current_index = 0
        case "K":
            if m.current_index == 0 {
                m.current_index = len(m.vault.Data.Folders) - 1
            } else {
                m.current_index -= 1
            }
            m.token_current_index = 0
        case "j":
            m.token_current_index = (m.token_current_index + 1) % len(m.vault.Data.Folders[m.current_index].Uris)
        case "k":
            if m.token_current_index == 0 {
                m.token_current_index = len(m.vault.Data.Folders[m.current_index].Uris) - 1
            } else {
                m.token_current_index -= 1
            }
        case "A":
            manager.PushScreen(InitializeNewFolderModel(&m.vault))
        case "e":
            manager.PushScreen(InitializeEditTokenModel(&m.vault, m.current_index, m.token_current_index))
        case "x":
            manager.PushScreen(InitializeDeleteTokenModel(&m.vault, m.current_index, m.token_current_index))
        case "ctrl+down":
            m.token_current_index += m.vault.MoveDown(m.current_index, m.token_current_index)
        case "ctrl+up":
            m.token_current_index -= m.vault.MoveUp(m.current_index, m.token_current_index)
        case "r":
            manager.PushScreen(InitializeMoveTokenModel(&m.vault, m.current_index, m.token_current_index))
        case "a":
            manager.PushScreen(InitializeAddTokenModel(&m.vault, m.current_index))
        }
    }

	return m, nil
}

// View
func (m DashboardModel) View() string {
    width, height, _ := term.GetSize(0)

    style := lipgloss.NewStyle().Height(height).Width(30).Padding(1, 3)
    folder_style := lipgloss.NewStyle().Width(30).Padding(1, 3)

    // Folders
    folders := make([]string, 0)

    for index, folder := range m.vault.Data.Folders {
        render_fn := folder_style.Render

        ui := lipgloss.JoinVertical(
            lipgloss.Left,
            m.styles.title.Render(folder.Name),
            m.styles.dimmed.Render(fmt.Sprintf("%d tokens", len(folder.Uris))),
        )

        if index == m.current_index {
            render_fn = folder_style.Copy().Background(lipgloss.Color("#1E1E2E")).
                Width(23).
                Padding(1, 2).
                BorderBackground(lipgloss.Color("#1E1E2E")).
                Border(lipgloss.ThickBorder(), false, false, false, true).Render
        }

        folders = append(folders, render_fn(ui))
    }

    // Tokens
    tokens := make([]string, 0)

    for index, uri := range m.vault.Data.Folders[m.current_index].Uris {
        style := lipgloss.NewStyle().
            Width(width - 30 - 2).
            Padding(1, 3).
            MarginBottom(1)
        title := lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
        issuer := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

        if index == m.token_current_index {
            style = style.Background(lipgloss.Color("#1E1E2E"))
            title = title.Background(lipgloss.Color("#1E1E2E")).Bold(true)
            issuer = issuer.Background(lipgloss.Color("#1E1E2E")).Bold(true)
        }

        _totp, _ := otp.NewKeyFromURL(uri)
        code, _ := totp.GenerateCode(_totp.Secret(), time.Now())

        icon := ""

        icon_spec, exists := m.vendor.Icons[_totp.Issuer()]

        if exists {
            style := lipgloss.NewStyle().
                Padding(1, 3).
                Foreground(lipgloss.Color("#FFFFFF")).
                Background(lipgloss.Color(fmt.Sprintf("#%s", icon_spec.Hex)));

            icon = style.Render(icon_spec.Unicode)
        }

        info := fmt.Sprintf("%s • %s", title.Render(_totp.AccountName()), issuer.Render(_totp.Issuer()))
        code_ui := make([]string, _totp.Digits())

        for _, code_char := range code {
            code_digit, _ := strconv.Atoi(string(code_char))

            code_ui = append(code_ui, strings.Trim(DIGIT_LIST[code_digit], "\n"))
        }

        spacing := style.GetWidth() - lipgloss.Width(info) - lipgloss.Width(lipgloss.JoinHorizontal(lipgloss.Center, code_ui...)) - 10 - lipgloss.Width(icon)

        tokens = append(tokens, style.Render(lipgloss.JoinHorizontal(lipgloss.Center, icon, title.Render("    "), info, strings.Repeat(" ", spacing), lipgloss.JoinHorizontal(lipgloss.Center, code_ui...))))
    }

    ui := []string {
        style.Render(lipgloss.JoinVertical(lipgloss.Left, folders...)),
    }

    if len(tokens) != 0 {
        ui = append(ui, lipgloss.JoinVertical(lipgloss.Left, tokens...))
    } else {
        placeholder := lipgloss.NewStyle().
            Width(width - 30 - 2).
            Height(height).
            Align(lipgloss.Center, lipgloss.Center)

        ui = append(ui, placeholder.Render("Press a to add some tokens"))
    }

    return lipgloss.JoinHorizontal(
        lipgloss.Left,
        ui...
    )
}

