package tokens

import (
	"errors"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	tlockform "github.com/eklairs/tlock/tlock-internal/form"
	"github.com/eklairs/tlock/tlock-internal/utils"
	tlockvault "github.com/eklairs/tlock/tlock-vault"
	tlockstyles "github.com/eklairs/tlock/tlock/styles"
	"github.com/pquerna/otp"
	"golang.org/x/term"
)

// Adds a validator to input to only accept int
func onlyInt(textinput textinput.Model) textinput.Model {
	return utils.ValidatorInteger(textinput)
}

// Validator for period
func periodValidator(_ *tlockvault.Vault, period string) error {
	if num, err := strconv.ParseInt(period, 10, 64); err == nil {
		if num < 1 {
			return errors.New("Period cannot be < 1")
		}
	}

	return nil
}

// Validator for digits
func digitValidator(_ *tlockvault.Vault, digits string) error {
	if num, err := strconv.ParseInt(digits, 10, 64); err == nil {
		if num < 1 {
			return errors.New("Digits cannot be < 1")
		}

		if num > 10 {
			return errors.New("Digits cannot be > 10")
		}
	}

	return nil
}

// Secret validator
func secretValidator(vault *tlockvault.Vault, secret string) error {
	// Validate
	_, err := vault.ValidateToken(secret)

	// Return
	return err
}

// Returns the form
func BuildForm(values map[string]string) tlockform.Form {
	// Initialize form
	form := tlockform.New()

	// Sets the value of the input box
	v := func(input textinput.Model, key string) textinput.Model {
		if value, ok := values[key]; ok {
			input.SetValue(value)
		}

		return input
	}

	// Add items
	form.AddInput("account", "Account Name", "Name of the account, like John Doe", v(components.InitializeInputBox("Account name goes here..."), "account"), []tlockform.Validator{})
	form.AddInput("issuer", "Issuer", "Name of the issuer, like GitHub", v(components.InitializeInputBox("Issuer name goes here..."), "issuer"), []tlockform.Validator{})
	form.AddInput("secret", "Secret", "The secret provided by the issuer", v(components.InitializeInputBox("The secret goes here..."), "secret"), []tlockform.Validator{secretValidator})
	form.AddOption("type", "Type", "Type of the token", []string{"TOTP", "HOTP"})
	form.AddOption("hash", "Hash", "Hashing algorithm for the token", []string{"SHA1", "SHA256", "SHA512"})
	form.AddInput("period", "Period", "Time to refresh the token", v(onlyInt(components.InitializeInputBoxCustomWidth("Time in seconds...", 24)), "period"), []tlockform.Validator{periodValidator})
	form.AddInput("counter", "Initial counter", "Initial counter for HOTP token", v(onlyInt(components.InitializeInputBoxCustomWidth("Initial counter...", 24)), "counter"), []tlockform.Validator{})
	form.AddInput("digits", "Digits", "Number of digits", v(onlyInt(components.InitializeInputBoxCustomWidth("Number of digits goes here...", 24)), "digits"), []tlockform.Validator{digitValidator})

	// Set default values
	form.Default = map[string]string{
		"account": "",
		"issuer":  "",
		"period":  "30",
		"counter": "0",
		"digits":  "6",
	}

	// Disable the counter box
	form.Disable("counter")

	// Run post init hook
	form.PostInit()

	// Return
	return form
}

// Renders the form
func RenderForm(ascii, description string, form tlockform.Form) string {
	// Items
	items := []string{
		tlockstyles.Styles.Title.Render(ascii), "",
		tlockstyles.Styles.SubText.Render("Add a new token [secret is required, rest are optional]"), "",
		form.Items[0].FormItem.View(), // Account name input
		form.Items[1].FormItem.View(), // Issuer name input
		form.Items[2].FormItem.View(), // Secret name input
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			form.Items[3].FormItem.View(), "   ",
			form.Items[4].FormItem.View(),
		), "",
	}

	// Render the input boxes based on choosen type
	inputGroup := lipgloss.JoinHorizontal(
		lipgloss.Left,
		form.Items[5].FormItem.View(), "   ",
		form.Items[7].FormItem.View(),
	)

	if form.Items[3].FormItem.Value() == "HOTP" {
		inputGroup = lipgloss.JoinHorizontal(
			lipgloss.Left,
			form.Items[6].FormItem.View(), "   ",
			form.Items[7].FormItem.View(),
		)
	}

	// Add the help menu
	items = append(items, inputGroup, "", tlockstyles.Help.View(addTokenKeys))

	// Return
	return lipgloss.JoinVertical(lipgloss.Center, items...)
}

// Creates a new viewport with the form rendered
func IntoViewport(ascii, description string, form tlockform.Form) viewport.Model {
	// Get size
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))

	// Rendered content
	content := RenderForm(ascii, description, form)

	// Initialize
	viewport := utils.DisableViewportKeys(viewport.New(85, min(height, lipgloss.Height(content))))
	viewport.SetContent(content)

	// Return
	return viewport
}

// Updates the viewport with the new form content
func UpdateViewport(ascii, description string, viewport *viewport.Model, form tlockform.Form) {
	// Get size
	_, height, _ := term.GetSize(int(os.Stdout.Fd()))

	// Rendered content
	content := RenderForm(ascii, description, form)

	// Update
	viewport.Height = min(height, lipgloss.Height(content))
	viewport.SetContent(content)
}

// Disables the form item based on the selected type
func DisableBasedOnType(form *tlockform.Form) {
	// Enable / Disable items based on the choosen type
	if form.Items[3].FormItem.Value() == "TOTP" {
		form.Disable("counter")
		form.Enable("period")
	}

	if form.Items[3].FormItem.Value() == "HOTP" {
		form.Enable("counter")
		form.Disable("period")
	}
}

// Converts string to token type
func toTokenType(tokentype string) tlockvault.TokenType {
	if tokentype == "HOTP" {
		return tlockvault.TokenTypeHOTP
	}

	return tlockvault.TokenTypeTOTP
}

// Converts string to opt algorithm
func toOtpAlgorithm(algo string) otp.Algorithm {
	switch algo {
	case "SHA256":
		return otp.AlgorithmSHA256
	case "SHA512":
		return otp.AlgorithmSHA512
	default:
		return otp.AlgorithmSHA1
	}
}

// Create a token from form data
func TokenFromFormData(data map[string]string) tlockvault.Token {
	return tlockvault.Token{
		Issuer:           data["issuer"],
		Account:          data["account"],
		Secret:           data["secret"],
		Type:             toTokenType(data["type"]),
		InitialCounter:   utils.ToInt(data["counter"]),
		Period:           utils.ToInt(data["period"]),
		Digits:           utils.ToInt(data["digits"]),
		HashingAlgorithm: toOtpAlgorithm(data["hash"]),
	}
}
