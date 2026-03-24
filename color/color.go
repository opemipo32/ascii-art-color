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

// function to color the gien text based on the color name passed as argument in cmd line
func Colorize(text, colorName string) (string, error) {
	code, err := Resolve(colorName)
	if err != nil {
		return "", err
	}
	return code + text + Reset, nil
}

// Getting the code for the colorName ( can be in RGB, HSL or ANSI)
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

// This code checks for RGB color.numbers validity
func validRGB(r, g, b int) bool {
	return r >= 0 && r <= 255 && g >= 0 && g <= 255 && b >= 0 && b <= 255
}

// This code converts HSL color → RGB color.
func hslToRGB(h, s, l int) (int, int, int) {
	hf := float64(h) / 360.0 // hue
	sf := float64(s) / 100.0 // saturation
	lf := float64(l) / 100.0 // lighting

	var r, g, b float64
	// if saturation is 0 means no color intensity
	// it becomes grey
	if sf == 0 {
		r, g, b = lf, lf, lf
	} else {
		var q float64 // upper brightness bound
		if lf < 0.5 {
			// if lightning is less than 50%
			q = lf * (1 + sf)
		} else {
			// if the lighting is equal to or greater than 50%
			q = lf + sf - lf*sf
		}

		p := 2*lf - q                // lower brightness bound
		r = hueToRGB(p, q, hf+1.0/3) // red
		g = hueToRGB(p, q, hf)       // green
		b = hueToRGB(p, q, hf-1.0/3) // blue
	}
	return int(r * 255), int(g * 255), int(b * 255)
}

// Moves upper and lower brightness according to
// the amount of hue
func hueToRGB(p, q, t float64) float64 {

	// p → minimum brightness
	// q → maximum brightness
	// t → position of the channel on the color wheel (normalized 0–1)

	// Hue is a circle, so t loops around.
	// Example: t = -0.2 → 0.8
	// Keeps the calculation inside 0–1.
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	// mapping t to brightness
	switch {
	case t < 1.0/6:
		// if  0 ≤ t < 1/6  then Ramp up from p → q
		return p + (q-p)*6*t
	case t < 1.0/2:
		// if 1/6 ≤ t < 1/2 then Stay at max q
		return q
	case t < 2.0/3:
		// if 1/2 ≤ t < 2/3 then Ramp down from q → p
		return p + (q-p)*(2.0/3-t)*6
	default:
		// if t ≥ 2/3 then Stay at min p
		return p
	}
}
