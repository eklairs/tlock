package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Type of simple icons index
type SimpleIcons struct {
    Icons []struct {
        Title string
        Hex string
    }
}

// Type of tlock's vendor
type TLockVendor struct {
    Icons map[string]struct {
        Unicode string
        Hex string
    }
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
    os.RemoveAll("vendor")
    os.RemoveAll("build")

    // Recreate
    os.Mkdir("vendor", os.ModePerm)
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

    vendor := TLockVendor{
        Icons: make(map[string]struct{Unicode string; Hex string}),
    }

    // Go go go
    for _, icon := range icons.Icons {
        value, ok := nerd_fonts[fmt.Sprintf("fa-%s", strings.ToLower(icon.Title))]

        if ok {
            vendor.Icons[icon.Title] = struct{Unicode string; Hex string}{ Unicode: value["char"], Hex: icon.Hex }
        }
    }

    // Write vendor
    dump, _ := json.Marshal(vendor);
    file, _ := os.Create(path.Join("vendor", "icons.json"))

    file.Write(dump)
}

