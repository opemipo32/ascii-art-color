package main

import (
	"strings"
	"testing"
)

func makeChars() [][]string {
	chars := make([][]string, 95)
	for i := range chars {
		row := strings.Repeat("-", 4)
		chars[i] = []string{row, row, row, row, row, row, row, row}
	}
	return chars
}

func TestPrintColoredRow_NoTarget_ColorsAll(t *testing.T) {
	chars := makeChars()
	printColoredRow("AB", 0, chars, "red", "")
}

func TestPrintColoredRow_WithTarget(t *testing.T) {
	chars := makeChars()
	printColoredRow("hello", 0, chars, "red", "ell")
}

func TestPrintColoredRow_TargetNotFound(t *testing.T) {
	chars := makeChars()
	printColoredRow("hello", 0, chars, "red", "xyz")
}
