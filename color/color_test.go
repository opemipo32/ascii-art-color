package color_test

import (
	"strings"
	"testing"

	"ascii-art/color"
)

func TestResolve_NamedColor(t *testing.T) {
	code, err := color.Resolve("red")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(code, "\033[") {
		t.Errorf("expected ANSI escape, got %q", code)
	}
}

func TestResolve_RGB(t *testing.T) {
	code, err := color.Resolve("rgb(255,0,0)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != "\033[38;2;255;0;0m" {
		t.Errorf("unexpected RGB code: %q", code)
	}
}

func TestResolve_HSL(t *testing.T) {

	code, err := color.Resolve("hsl(0,100,50)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(code, "\033[38;2;") {
		t.Errorf("expected 24-bit ANSI code, got %q", code)
	}
}

func TestResolve_Unsupported(t *testing.T) {
	_, err := color.Resolve("ultraviolet")
	if err == nil {
		t.Fatal("expected error for unsupported color, got nil")
	}
}

func TestColorize(t *testing.T) {
	out, err := color.Colorize("hello", "red")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "hello") {
		t.Error("colorized output should contain original text")
	}
	if !strings.HasSuffix(out, color.Reset) {
		t.Error("colorized output should end with reset code")
	}
}
