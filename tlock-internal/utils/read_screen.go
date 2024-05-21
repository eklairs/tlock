package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/kbinani/screenshot"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

// Returns if the current session is wayland
func isWayland() bool {
	return os.Getenv("XDG_SESSION_TYPE") == "wayland"
}

// Error when no token on the screen is found
var TOKEN_NOT_FOUND_ERR = errors.New("No token found on the screen")

// Reads QRCode from screen and returns the found data (not for wayland)
func ReadTokenFromScreenNoWayland() (*string, error) {
	// Capture rect
	image, err := screenshot.CaptureRect(screenshot.GetDisplayBounds(0))
	if err != nil {
		return nil, err
	}

	// Create bitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(image)
	if err != nil {
		return nil, err
	}

	// Create reader
	qrReader := qrcode.NewQRCodeReader()

	// Decode
	if result, err := qrReader.Decode(bmp, nil); err == nil {
		// Get fetched URI
		uri := result.String()

		// Return
		return &uri, nil
	}

	// No token found
	return nil, TOKEN_NOT_FOUND_ERR
}

// Reads QRCode form screen and returns the found data (only for wayland)
func ReadTokenFromScreenWaylandOnly() (*string, error) {
	// Make grim command
	_, err := exec.Command("grim").Output()

	// Check for error
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot run grim command: %s", err))
	}

	return nil, nil
}

// Reads QRCode from screen based on the session type
func ReadTokenFromScreen() (*string, error) {
	// Check if it is wayland
	if isWayland() {
		return ReadTokenFromScreenWaylandOnly()
	}

	// Use normal function
	return ReadTokenFromScreenNoWayland()
}
