package main

import (
	"fmt"
	"io"
	"time"
)

// SwapRecord represents a single swap event for the swap monitor panel.
type SwapRecord struct {
	PID       string
	From      string // tier name: "VRAM", "RAM", "SSD"
	To        string // tier name: "VRAM", "RAM", "SSD"
	Timestamp time.Time
	SizeBytes uint64
}

// SwapReader supplies recent swap events for read-only display.
type SwapReader interface {
	SwapEvents() []SwapRecord
}

// writeSwapPanel renders the swap monitor panel (Panel 4).
// Shows recent tier movements: VRAM → RAM → SSD.
func writeSwapPanel(w io.Writer, swaps SwapReader) {
	var lines []string

	var events []SwapRecord
	if swaps != nil {
		events = swaps.SwapEvents()
	}

	if len(events) == 0 {
		lines = append(lines, dim("(no swap activity)"))
		writeBox(w, "Swap Monitor", boxWidth, lines)
		return
	}

	// Header.
	lines = append(lines, fmt.Sprintf("%-8s %-6s %-4s %-6s %s",
		bold("PID"), bold("FROM"), "→", bold("TO"), bold("SIZE")))

	// Show up to 8 most recent events.
	start := 0
	if len(events) > 8 {
		start = len(events) - 8
	}
	for _, ev := range events[start:] {
		sizeStr := formatBytes(ev.SizeBytes)
		arrow := colorize("→", ansiYellow)
		lines = append(lines, fmt.Sprintf("%-8s %-6s %s %-6s %s",
			truncate(ev.PID, 8), ev.From, arrow, ev.To, sizeStr))
	}

	writeBox(w, "Swap Monitor", boxWidth, lines)
}

// formatBytes formats a byte count into a human-readable string.
func formatBytes(b uint64) string {
	const (
		mib = 1024 * 1024
		gib = 1024 * mib
	)
	switch {
	case b >= gib:
		return fmt.Sprintf("%.1f GiB", float64(b)/float64(gib))
	case b >= mib:
		return fmt.Sprintf("%.0f MiB", float64(b)/float64(mib))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
