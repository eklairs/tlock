package tokens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/eklairs/tlock/tlock-internal/modelmanager"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
	"github.com/pquerna/otp/totp"
)

var (
    editedAccountName = ""
    editedIssuerName = ""
    editedSecret = ""
    editedIsTotp = true
    editedDigits = 0
    editedHashFn = otp.AlgorithmSHA256
)

type EditTokenModel struct {
    form *huh.Form
    old string
}

func InitializeEditTokenModel(uri string) EditTokenModel {
    key, _ := otp.NewKeyFromURL(uri)

    editedAccountName = key.AccountName()
    editedIssuerName = key.Issuer()
    editedDigits = key.Digits().Length()
    editedHashFn = key.Algorithm()
    editedSecret = key.Secret()
    editedIsTotp = key.Type() == "totp"

    digitsStr := fmt.Sprintf("%d", editedDigits)

    form := huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("Account name").
                Value(&editedAccountName).
                Description("Name of the account, like John Doe"),
            huh.NewInput().
                Title("Issuer name").
                Value(&editedIssuerName).
                Description("Name of the issuer, like GitHub"),
            huh.NewInput().
                Title("Secret").
                Value(&editedSecret).
                Description("Enter the secret provided by the issuer"),
            huh.NewConfirm().
                Title("Type").
                Description("The type of the token").
                Value(&editedIsTotp).
                Negative("HOTP").
                Affirmative("TOTP"),
            huh.NewSelect[otp.Algorithm]().
                Title("Hash function").
                Description("The hash function for the token").
                Value(&editedHashFn).
                Options(
                    huh.NewOption("SHA-256", otp.AlgorithmSHA256),
                    huh.NewOption("SHA-512", otp.AlgorithmSHA512),
                    huh.NewOption("SHA-1", otp.AlgorithmSHA1),
                    huh.NewOption("MD5", otp.AlgorithmMD5),
                ),
            huh.NewInput().
                Title("Digits").
                Value(&digitsStr).
                Description("Number of digits"),
        ),
    )

    return EditTokenModel{
        form: form,
        old: uri,
    }
}

func (model EditTokenModel) Init() tea.Cmd {
    return model.form.Init()
}

func (model EditTokenModel) Update(msg tea.Msg, manager *modelmanager.ModelManager) (modelmanager.Screen, tea.Cmd) {
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
                Issuer: editedIssuerName,
                AccountName: editedAccountName,
                Period: 30,
                Secret: []byte(editedSecret),
                Algorithm: editedHashFn,
            })

            cmds = append(cmds, func() tea.Msg {
                return EditTokenMsg{
                    Old: model.old,
                    New: totpKey.URL(),
                }
            })
        } else {
            hotpKey, _ := hotp.Generate(hotp.GenerateOpts{
                Issuer: editedIssuerName,
                AccountName: editedAccountName,
                Secret: []byte(editedSecret),
                Algorithm: editedHashFn,
            })

            cmds = append(cmds, func() tea.Msg {
                return EditTokenMsg{
                    Old: model.old,
                    New: hotpKey.URL(),
                }
            })
        }

        manager.PopScreen()
    }

    return model, tea.Batch(cmds...)
}

func (model EditTokenModel) View() string {
    return model.form.View()
}

