package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"ascii-art/banner"
	"ascii-art/color"
)

const usage = "Usage: go run . [OPTION] [STRING]\nEX: go run . --color=<color> <substring to be colored> \"something\""

func main() {
	if len(os.Args) == 0 {
		fmt.Println(usage)
		return
	}

	if len(os.Args) != 3 && len(os.Args) != 5 && len(os.Args) != 2 && len(os.Args) != 4 {
		fmt.Println("length Issue..")
		return
	}

	// ========== Extracting SubString and Input ======
	colorTarget := ""
	bannerName := "standard"
	var input string

	option := ""

	switch len(os.Args) {
	case 4:
		colorTarget = os.Args[2]
		input = os.Args[3]
		option = os.Args[1]

	case 5:
		colorTarget = os.Args[2]
		input = os.Args[3]
		bannerName = os.Args[4]
		option = os.Args[1]

	case 3:
		input = os.Args[2]
		option = os.Args[1]
	case 2:
		input = os.Args[1]
	}

	input = strings.ReplaceAll(input, "\\n", "\n")

	if !strings.HasPrefix(option, "--color=") {
		if option != "" {
			fmt.Println(usage)
			return
		}
	}

	if !isValidBanner(bannerName) {
		fmt.Println(usage)
		return
	}

	colorValue := strings.ToLower(strings.TrimPrefix(option, "--color="))
	if _, err := color.Resolve(colorValue); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	chars, err := banner.Load(bannerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading banner: %v\n", err)
		os.Exit(1)
	}

	// Printing A new line for empty input
	if input == "" {
		fmt.Println()
		return
	}

	// Escaping all newline from the input
	input = strings.ReplaceAll(input, `\n`, "\n")
	// spliting the input into deffrent lines
	lines := strings.Split(input, "\n")

	// looping through each line
	for i, line := range lines {
		// printing new line for empty line found
		if line == "" {
			if i < len(lines)-1 {
				fmt.Println()
			}
			continue
		}
		for _, ch := range line {
			idx := int(ch) - 32
			if idx < 0 || idx > 94 {
				fmt.Fprintf(os.Stderr, "Error: character '%c' (code %d) is not supported\n", ch, ch)
				os.Exit(1)
			}
		}

		for row := 0; row < 8; row++ {
			if colorValue == "" {
				for _, ch := range line {
					fmt.Print(chars[int(ch)-32][row])
				}
			} else {
				printColoredRow(line, row, chars, colorValue, colorTarget)
			}
			fmt.Println()
		}
	}
}
func printColoredRow(line string, row int, chars [][]string, colorValue, colorTarget string) {
	runes := []rune(line)
	colored := make([]bool, len(runes))

	if colorTarget == "" {
		for i := range colored {
			// setting the color target to true for the
			// whole if no sub string is specified
			colored[i] = true
		}
	} else {
		// conerting the substring into a slice of runes
		targetRunes := []rune(colorTarget)
		tLen := len(targetRunes)
		for i := 0; i <= len(runes)-tLen; i++ {
			match := true
			for j := 0; j < tLen; j++ {
				if runes[i+j] != targetRunes[j] {
					match = false
					break
				}
			}
			if match {
				for j := 0; j < tLen; j++ {
					colored[i+j] = true
				}
			}
		}
	}

	// Checking if the colorValue is valid
	ansiCode, _ := color.Resolve(colorValue)
	for i, ch := range runes {
		art := chars[int(ch)-32][row]
		if colored[i] {
			fmt.Print(ansiCode + art + color.Reset)
		} else {
			fmt.Print(art)
		}
	}
}

// used to check if the
// input banner name
// is one we work with
func isValidBanner(name string) bool {
	validBanners := []string{"standard", "shadow", "thinkertoy"}

	return slices.Contains(validBanners, name)
}
