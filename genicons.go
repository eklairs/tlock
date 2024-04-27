package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

var OUT_DIR = path.Join("tlock-vendor", "icons")

// Priority list
var PRIORITY = []string{
	"seti",
	"fae",
	"md",
}

// Type of simple icons index
type SimpleIcons struct {
	Icons []struct {
		Title string
		Hex   string
	}
}

type Icon struct {
	Unicode string
	Hex     string
}

// Vendor icons
type VendorIcons struct {
	Icons map[string]Icon
}

func parse_simple_icons_map() SimpleIcons {
	raw, err := os.ReadFile(path.Join("build", "simple-icons", "_data", "simple-icons.json"))

	if err != nil {
		panic("Failed to read icon map")
	}

	out := SimpleIcons{}

	if err = json.Unmarshal(raw, &out); err != nil {
		panic("Failed to parse icon map")
	}

	return out
}

func parse_nerd_fonts_map() map[string]map[string]string {
	raw, err := os.ReadFile(path.Join("build", "nerdfonts.json"))

	if err != nil {
		panic("Failed to read icon map")
	}

	out := make(map[string]map[string]string)

	if err = json.Unmarshal(raw, &out); err != nil {
		panic("Failed to parse icon map")
	}

	return out
}

func main() {
	// Clear current vendor and build dirs
	os.RemoveAll(OUT_DIR)
	os.RemoveAll("build")

	// Make directory
	os.MkdirAll(OUT_DIR, os.ModePerm)
	os.Mkdir("build", os.ModePerm)

	// Clone simple icons
	fmt.Println("=> Cloning simple-icons git repository...")
	exec.Command("git", "clone", "https://github.com/simple-icons/simple-icons.git", path.Join("build", "simple-icons"), "--depth", "1").Output()

	// Download nerd fonts glyph names
	fmt.Println("=> Downloading nerd fonts glyph names")
	exec.Command("curl", "https://raw.githubusercontent.com/ryanoasis/nerd-fonts/master/glyphnames.json", "-o", path.Join("build", "nerdfonts.json")).Output()
	nerd_fonts := parse_nerd_fonts_map()

	// Parse icons map
	fmt.Println("=> Parsing simple-icons' icons map")
	icons := parse_simple_icons_map()

	vendor := VendorIcons{
		Icons: make(map[string]Icon),
	}

	// Go go go
	for _, icon := range icons.Icons {
		// Fetch from the list based on priority
	priority:
		for _, provider := range PRIORITY {
			// Fetch from the provider
			value, ok := nerd_fonts[fmt.Sprintf("%s-%s", provider, strings.ToLower(icon.Title))]

			// If found, add
			if ok {
				// Add
				vendor.Icons[icon.Title] = Icon{Unicode: value["char"], Hex: icon.Hex}

				// Break
				break priority
			}
		}
	}

	// Write vendor
	dump, _ := json.Marshal(vendor)
	file, _ := os.Create(path.Join(OUT_DIR, "icons.json"))

	file.Write(dump)
}
