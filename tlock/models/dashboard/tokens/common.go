package tokens

import (
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/eklairs/tlock/tlock-internal/components"
	"github.com/eklairs/tlock/tlock-internal/form"
	tlockstyles "github.com/eklairs/tlock/tlock/styles"
	"golang.org/x/term"
)

// Returns the form
func BuildForm() form.Form {
    // Initialize form
    form := form.New();

    // Add items
    form.AddInput("account", "Account Name", "Name of the account, like John Doe", components.InitializeInputBox("Account name goes here..."), []textinput.ValidateFunc{})
    form.AddInput("issuer", "Issuer", "Name of the issuer, like GitHub", components.InitializeInputBox("Issuer name goes here..."), []textinput.ValidateFunc{})
    form.AddInput("secret", "Secret", "The secret provided by the issuer", components.InitializeInputBox("The secret goes here..."), []textinput.ValidateFunc{})
    form.AddOption("type", "Type", "Type of the token", []string{"TOTP", "HOTP"})
    form.AddOption("hash", "Hash", "Hashing algorithm for the token", []string{"SHA1", "SHA256", "SHA512"})
    form.AddInput("period", "Period", "Time to refresh the token", components.InitializeInputBoxCustomWidth("Time in seconds...", 24), []textinput.ValidateFunc{})
    form.AddInput("counter", "Initial counter", "Initial counter for HOTP token", components.InitializeInputBoxCustomWidth("Initial counter...", 24), []textinput.ValidateFunc{})
    form.AddInput("digits", "Digits", "Number of digits", components.InitializeInputBoxCustomWidth("Number of digits goes here...", 24), []textinput.ValidateFunc{})

    // Disable the counter box
    form.Disable("counter")

    // Run post init hook
    form.PostInit()

    // Return
    return form
}

// Renders the form
func RenderForm(ascii, description string, form form.Form) string {
    // Items
	items := []string{
		tlockstyles.Styles.Title.Render(ascii), "",
		tlockstyles.Styles.SubText.Render("Add a new token [secret is required, rest are optional]"), "",
		form.Items[0].FormItem.View(), // Account name input
		form.Items[1].FormItem.View(), // Account name input
		form.Items[2].FormItem.View(), // Account name input
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
func IntoViewport(ascii, description string, form form.Form) viewport.Model {
    // Get size
    _, height, _ := term.GetSize(int(os.Stdout.Fd()))

    // Rendered content
    content := RenderForm(ascii, description, form)

    // Initialize
    viewport := viewport.New(85, min(height, lipgloss.Height(content)))
    viewport.SetContent(content)

    // Return
    return viewport
}

// Updates the viewport with the new form content
func UpdateViewport(ascii, description string, viewport *viewport.Model, form form.Form) {
    // Get size
    _, height, _ := term.GetSize(int(os.Stdout.Fd()))

    // Rendered content
    content := RenderForm(ascii, description, form)

    // Update
    viewport.Height = min(height, lipgloss.Height(content))
    viewport.SetContent(content)
}

// Disables the form item based on the selected type
func DisableBasedOnType(form *form.Form) {
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
