package main

import (
	"fmt"
	"os"
	"strings"

	"ascii-art/banner"
	"ascii-art/color"
)

const usage = "Usage: go run . [OPTION] [STRING]\nEX: go run . --color=<color> <substring to be colored> \"something\""

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		return
	}

	colorValue := ""
	colorTarget := ""
	var remaining []string

	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--color=") {
			colorValue = strings.TrimPrefix(args[i], "--color=")
			if colorValue == "" {
				fmt.Fprintln(os.Stderr, usage)
				os.Exit(1)
			}
		} else if strings.HasPrefix(args[i], "--") {
			fmt.Fprintln(os.Stderr, usage)
			os.Exit(1)
		} else {
			remaining = append(remaining, args[i])
		}
	}

	input := ""
	bannerName := "standard"

	if colorValue != "" {
		if len(remaining) == 0 || len(remaining) > 3 {
			fmt.Fprintln(os.Stderr, usage)
			os.Exit(1)
		}

		if len(remaining) == 1 {
			colorTarget = ""
			input = remaining[0]
		} else if len(remaining) == 2 {
			if isValidBanner(remaining[1]) {
				colorTarget = ""
				input = remaining[0]
				bannerName = remaining[1]
			} else {
				colorTarget = remaining[0]
				input = remaining[1]
			}
		} else if len(remaining) == 3 {
			colorTarget = remaining[0]
			input = remaining[1]
			bannerName = remaining[2]
		}
	} else {
		if len(remaining) < 1 || len(remaining) > 2 {
			fmt.Fprintln(os.Stderr, usage)
			os.Exit(1)
		}
		input = remaining[0]
		if len(remaining) == 2 {
			bannerName = remaining[1]
		}
	}
	if colorValue != "" {
		if _, err := color.Resolve(colorValue); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	chars, err := banner.Load(bannerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading banner: %v\n", err)
		os.Exit(1)
	}

	if input == "" {
		fmt.Println()
		return
	}

	input = strings.ReplaceAll(input, `\n`, "\n")
	lines := strings.Split(input, "\n")

	for i, line := range lines {
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
			colored[i] = true
		}
	} else {
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

func isValidBanner(name string) bool {
	validBanners := []string{"standard", "shadow", "thinkertoy"}
	for _, b := range validBanners {
		if name == b {
			return true
		}
	}
	return false
}
