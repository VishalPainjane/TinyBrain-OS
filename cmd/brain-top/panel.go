package main

import (
	"fmt"
	"io"
	"strings"
)

// boxWidth is the default panel width for box-drawing frames.
const boxWidth = 60

// writeBox draws a Unicode box frame around content lines with a title.
//
//	┌─ Title ──────────────────────────────┐
//	│ content line 1                       │
//	│ content line 2                       │
//	└──────────────────────────────────────┘
func writeBox(w io.Writer, title string, width int, lines []string) {
	if width <= 0 {
		width = boxWidth
	}
	innerWidth := width - 2 // subtract left and right border characters

	// Top border: ┌─ Title ─────┐
	titleLabel := ""
	if title != "" {
		titleLabel = " " + bold(title) + " "
	}
	titleVisualLen := 0
	if title != "" {
		titleVisualLen = len(title) + 2 // space + title + space (excluding ANSI codes)
	}
	topDashes := innerWidth - titleVisualLen - 1 // -1 for the leading dash
	if topDashes < 0 {
		topDashes = 0
	}
	fmt.Fprintf(w, "┌─%s%s┐\n", titleLabel, strings.Repeat("─", topDashes))

	// Content lines.
	for _, line := range lines {
		visualLen := visualLength(line)
		padding := innerWidth - visualLen
		if padding < 0 {
			padding = 0
		}
		fmt.Fprintf(w, "│ %s%s│\n", line, strings.Repeat(" ", padding))
	}

	// If no content lines, draw an empty line.
	if len(lines) == 0 {
		fmt.Fprintf(w, "│%s│\n", strings.Repeat(" ", innerWidth+1))
	}

	// Bottom border: └─────────────┘
	fmt.Fprintf(w, "└%s┘\n", strings.Repeat("─", innerWidth+1))
}

// visualLength returns the display width of a string, excluding ANSI escape sequences.
func visualLength(s string) int {
	inEscape := false
	length := 0
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEscape = false
			}
			continue
		}
		length++
	}
	return length
}

// padRight pads s to the given width with spaces.
// If s is already wider than width, it is returned unchanged.
func padRight(s string, width int) string {
	vl := visualLength(s)
	if vl >= width {
		return s
	}
	return s + strings.Repeat(" ", width-vl)
}
