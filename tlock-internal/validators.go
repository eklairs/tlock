package tlockinternal

import (
	"errors"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
)

// Adds a validator function to the given input box
func Validator(input textinput.Model, function textinput.ValidateFunc) textinput.Model {
	input.Validate = function

	return input
}

// Validator for checking the value is an integer
func ValidatorInteger(input textinput.Model) textinput.Model {
	return Validator(input, func(inputValue string) error {
		_, err := strconv.ParseInt(inputValue, 10, 64)

		return err
	})
}

// Validator for checking the value is an integer and doesnt not allow zero
func ValidatorIntegerNo0(input textinput.Model) textinput.Model {
	return Validator(input, func(inputValue string) error {
		num, err := strconv.ParseInt(inputValue, 10, 64)

		if num == 0 {
			return errors.New("Value cannot be zero")
		}

		return err
	})
}
