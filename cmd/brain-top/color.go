package main

import "fmt"

// ANSI escape codes for terminal styling.
// Supports any terminal with ANSI support (Windows Terminal 1511+, all Linux/macOS).
const (
	ansiReset = "\033[0m"
	ansiBold  = "\033[1m"
	ansiDim   = "\033[2m"

	ansiRed     = "\033[31m"
	ansiGreen   = "\033[32m"
	ansiYellow  = "\033[33m"
	ansiBlue    = "\033[34m"
	ansiMagenta = "\033[35m"
	ansiCyan    = "\033[36m"
	ansiWhite   = "\033[37m"

	ansiGreenBg = "\033[42m"
	ansiRedBg   = "\033[41m"
)

// colorize wraps text with an ANSI escape code and reset.
func colorize(text, code string) string {
	if !colorEnabled {
		return text
	}
	return code + text + ansiReset
}

// bold wraps text in bold ANSI styling.
func bold(text string) string {
	return colorize(text, ansiBold)
}

// dim wraps text in dim ANSI styling.
func dim(text string) string {
	return colorize(text, ansiDim)
}

// colorEnabled controls whether ANSI colours are emitted.
// Disable by setting NO_COLOR=1 in the environment.
var colorEnabled = true

// stateColor returns the ANSI colour code for a process state string.
func stateColor(state string) string {
	switch state {
	case "RUNNING":
		return ansiGreen
	case "READY":
		return ansiCyan
	case "WAITING":
		return ansiYellow
	case "PREEMPTED":
		return ansiRed
	case "HIBERNATED":
		return ansiBlue
	case "NEW":
		return ansiWhite
	case "TERMINATED":
		return ansiDim
	default:
		return ansiWhite
	}
}

// bar renders a utilization bar like [████████░░░░] 67%.
// width is the number of characters inside the brackets.
func bar(used, total uint64, width int) string {
	if width <= 0 {
		width = 20
	}
	if total == 0 {
		empty := make([]byte, width)
		for i := range empty {
			empty[i] = ' '
		}
		return fmt.Sprintf("[%s]  0%%", string(empty))
	}

	pct := float64(used) / float64(total)
	if pct > 1.0 {
		pct = 1.0
	}
	filled := int(pct * float64(width))
	if filled > width {
		filled = width
	}

	fillChar := "█"
	emptyChar := "░"

	// Choose colour based on utilization percentage.
	barColor := ansiGreen
	switch {
	case pct >= 0.9:
		barColor = ansiRed
	case pct >= 0.7:
		barColor = ansiYellow
	}

	result := "["
	if colorEnabled {
		result += barColor
	}
	for i := 0; i < filled; i++ {
		result += fillChar
	}
	if colorEnabled {
		result += ansiReset + ansiDim
	}
	for i := filled; i < width; i++ {
		result += emptyChar
	}
	if colorEnabled {
		result += ansiReset
	}
	result += fmt.Sprintf("] %3.0f%%", pct*100)
	return result
}
