package tokens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
)

var (
	accountName = ""
	issuerName  = ""
	secret      = ""
	isTotp      = true
	digits      = 0
	hashFn      = otp.AlgorithmSHA256
)

type AddTokenModel struct {
	form *huh.Form
}

func InitializeAddTokenModel() AddTokenModel {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Account name").
				Value(&accountName).
				Description("Name of the account, like John Doe"),
			huh.NewInput().
				Title("Issuer name").
				Value(&issuerName).
				Description("Name of the issuer, like GitHub"),
			huh.NewInput().
				Title("Secret").
				Value(&secret).
				Description("Enter the secret provided by the issuer"),
			huh.NewConfirm().
				Title("Type").
				Description("The type of the token").
				Value(&isTotp).
				Negative("HOTP").
				Affirmative("TOTP"),
			huh.NewSelect[otp.Algorithm]().
				Title("Hash function").
				Description("The hash function for the token").
				Value(&hashFn).
				Options(
					huh.NewOption("SHA-256", otp.AlgorithmSHA256),
					huh.NewOption("SHA-512", otp.AlgorithmSHA512),
					huh.NewOption("SHA-1", otp.AlgorithmSHA1),
					huh.NewOption("MD5", otp.AlgorithmMD5),
				),
			huh.NewInput().
				Title("Digits").
				Description("Number of digits"),
		),
	)

	return AddTokenModel{
		form: form,
	}
}

func (model AddTokenModel) Init() tea.Cmd {
	return model.form.Init()
}

func (model AddTokenModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "esc":
			manager.PopScreen()
		}
	}

	// Change form sumbits
	_, cmd := model.form.Update(msg)
	cmds = append(cmds, cmd)

	// Form is sumbitted
	if model.form.State == huh.StateCompleted {
		if isTotp {
			totpKey, _ := totp.Generate(totp.GenerateOpts{
				Issuer:      issuerName,
				AccountName: accountName,
				Period:      30,
				Secret:      []byte(secret),
				Algorithm:   hashFn,
			})

			cmds = append(cmds, func() tea.Msg {
				return AddTokenMsg{
					URI: totpKey.URL(),
				}
			})
		} else {
			hotpKey, _ := hotp.Generate(hotp.GenerateOpts{
				Issuer:      issuerName,
				AccountName: accountName,
				Secret:      []byte(secret),
				Algorithm:   hashFn,
			})

			cmds = append(cmds, func() tea.Msg {
				return AddTokenMsg{
					URI: hotpKey.URL(),
				}
			})
		}

		manager.PopScreen()
	}

	return model, tea.Batch(cmds...)
}

func (model AddTokenModel) View() string {
	return model.form.View()
}
