package utils

import (
	"errors"
	"fmt"
	"image"
	"os"
	"os/exec"

	_ "image/png"

	"github.com/kbinani/screenshot"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

// Loads an image from file
func getImageFromFilePath(filePath string) (image.Image, error) {
	// Read file path
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	// Close on scope end
	defer f.Close()

	// Decode
	image, _, err := image.Decode(f)

	// Return
	return image, err
}

// Returns if the current session is wayland
func isWayland() bool {
	return os.Getenv("XDG_SESSION_TYPE") == "wayland"
}

// Reads the QRCode from the image
func readFromImage(image image.Image) (*string, error) {
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

// Error when no token on the screen is found
var TOKEN_NOT_FOUND_ERR = errors.New("No token found on the screen")

// Reads QRCode from screen and returns the found data (not for wayland)
func ReadTokenFromScreenNoWayland() (*string, error) {
	// Capture rect
	image, err := screenshot.CaptureRect(screenshot.GetDisplayBounds(0))
	if err != nil {
		return nil, err
	}

	// Return
	return readFromImage(image)
}

// Reads QRCode form screen and returns the found data (only for wayland)
func ReadTokenFromScreenWaylandOnly() (*string, error) {
	// Out path
	out := "/tmp/tlock_screenshot.png"

	// Make grim command
	_, err := exec.Command("grim", out).Output()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot run grim command: %s", err))
	}

	// Load screenshot
	image, err := getImageFromFilePath(out)
	if err != nil {
		return nil, err
	}

	// Read
	return readFromImage(image)
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
