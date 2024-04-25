package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/charmbracelet/lipgloss"
	tlockstyles "github.com/eklairs/tlock/tlock-styles"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

const LOG_PREFIX = "=>"

const BUILD_DIR = "build"

var VENDOR_DIR = path.Join("tlock-vendor", "proprietary")

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}

func parse_name(file string) string {
	name := strings.TrimSuffix(file, ".css")
	parts := Map(strings.Split(name, "_"), func(part string) string { return strings.ToUpper(string(part[0])) + part[1:] })

	return strings.Join(parts, " ")
}

// Quick log utility
func debug(text string) {
	// Styles
	prefix := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	message := lipgloss.NewStyle().Foreground(lipgloss.Color("4"))

	// Render them all
	fmt.Printf("%s %s\n", prefix.Render(LOG_PREFIX), message.Render(text))
}

func main() {
	// Make required directories
	os.MkdirAll(BUILD_DIR, os.ModePerm)
	os.MkdirAll(VENDOR_DIR, os.ModePerm)

	// Debug
	debug("Generating tlock proprietary vendor...")

	// Clone monkeytype
	debug("Cloning monkeytypegame/monkeytype...")
	exec.Command("git", "clone", "https://github.com/monkeytypegame/monkeytype", path.Join(BUILD_DIR, "monkeytype"), "--depth", "1").Output()

	// Iterate over all the theme css
	themes_dir := path.Join(BUILD_DIR, "monkeytype", "frontend", "static", "themes")
	css_entries, _ := os.ReadDir(themes_dir)

	// Themes list
	themes := make([]tlockstyles.Theme, 0)

	for _, file := range css_entries {
		if strings.HasSuffix(file.Name(), ".css") {
			// Read raw css
			css_raw, _ := os.ReadFile(path.Join(themes_dir, file.Name()))

			// Parse
			lexer := css.NewLexer(parse.NewInput(strings.NewReader(string(css_raw[:]))))

			// Theme instance
			theme := tlockstyles.Theme{
				Name: parse_name(file.Name()),
			}

			// Flag to specify to start collecting values for colors
			startCollectingVars := false
			previousCustomName := ""

			// Iter till we find :root
		parser:
			for {
				tt, text := lexer.Next()

				switch tt {
				case css.ErrorToken:
					if lexer.Err() == io.EOF {
						break parser
					}
				case css.IdentToken:
					if string(text[:]) == "root" {
						startCollectingVars = true
					}
				case css.CustomPropertyNameToken:
					if startCollectingVars {
						previousCustomName = string(text[:])
					}
				case css.HashToken:
					if startCollectingVars {
						hex_color := lipgloss.Color(string(text[:]))

						switch previousCustomName {
						case "--bg-color":
							theme.SubAlt = hex_color
						case "--main-color":
							theme.Accent = hex_color
						case "--sub-color":
							theme.Sub = hex_color
						case "--sub-alt-color":
							theme.Background = hex_color
						case "--text-color":
							theme.Text = hex_color
						case "--error-color":
							theme.Error = hex_color
						}
					}
				}
			}

			themes = append(themes, theme)
		}
	}

	debug("Writing themes...")

	// Dump
	themes_dump, _ := json.Marshal(themes)

	// Create and write
	file, _ := os.Create(path.Join(VENDOR_DIR, "themes.json"))

	// Write
	file.Write(themes_dump)

	// Debug
	debug("Done writing themes!")
}
