package color

import "fmt"

const Reset = "\033[0m"

var namedColors = map[string]string{
	"black":   "\033[30m",
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",
	"orange":  "\033[38;5;214m",
	"pink":    "\033[38;5;213m",
	"purple":  "\033[38;5;129m",
}

func Colorize(text, colorName string) (string, error) {
	code, err := Resolve(colorName)
	if err != nil {
		return "", err
	}
	return code + text + Reset, nil
}
func Resolve(colorName string) (string, error) {

	if code, ok := namedColors[colorName]; ok {
		return code, nil
	}

	var r, g, b int
	if n, _ := fmt.Sscanf(colorName, "rgb(%d,%d,%d)", &r, &g, &b); n == 3 {
		if validRGB(r, g, b) {
			return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b), nil
		}
		return "", fmt.Errorf("rgb values must be between 0 and 255")
	}

	var h, s, l int
	if n, _ := fmt.Sscanf(colorName, "hsl(%d,%d,%d)", &h, &s, &l); n == 3 {
		r, g, b := hslToRGB(h, s, l)
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b), nil
	}

	return "", fmt.Errorf("unsupported color: %q", colorName)
}

func validRGB(r, g, b int) bool {
	return r >= 0 && r <= 255 && g >= 0 && g <= 255 && b >= 0 && b <= 255
}

func hslToRGB(h, s, l int) (int, int, int) {
	hf := float64(h) / 360.0
	sf := float64(s) / 100.0
	lf := float64(l) / 100.0

	var r, g, b float64
	if sf == 0 {
		r, g, b = lf, lf, lf
	} else {
		var q float64
		if lf < 0.5 {
			q = lf * (1 + sf)
		} else {
			q = lf + sf - lf*sf
		}
		p := 2*lf - q
		r = hueToRGB(p, q, hf+1.0/3)
		g = hueToRGB(p, q, hf)
		b = hueToRGB(p, q, hf-1.0/3)
	}
	return int(r * 255), int(g * 255), int(b * 255)
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	switch {
	case t < 1.0/6:
		return p + (q-p)*6*t
	case t < 1.0/2:
		return q
	case t < 2.0/3:
		return p + (q-p)*(2.0/3-t)*6
	default:
		return p
	}
}
