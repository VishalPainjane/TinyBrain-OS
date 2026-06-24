package main

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/hardware"
	"github.com/VishalPainjane/TinyBrain-OS/internal/process"
	"github.com/VishalPainjane/TinyBrain-OS/internal/scheduler"
)

// ProcessReader supplies process table rows for read-only display.
type ProcessReader interface {
	List() []process.Process
}

// QueueReader supplies MLFQ queue depths for read-only display.
type QueueReader interface {
	Depths() [scheduler.NumLevels]int
}

// processTableReader adapts ProcessTable to ProcessReader.
type processTableReader struct {
	table *process.ProcessTable
}

func (r processTableReader) List() []process.Process {
	if r.table == nil {
		return nil
	}
	return r.table.List()
}

// mlfqQueueReader adapts MLFQScheduler to QueueReader.
type mlfqQueueReader struct {
	sched *scheduler.MLFQScheduler
}

func (r mlfqQueueReader) Depths() [scheduler.NumLevels]int {
	if r.sched == nil {
		return [scheduler.NumLevels]int{}
	}
	return r.sched.QueueDepths()
}

// renderSnapshot writes a full brain-top dashboard frame to w.
// See docs/architecture/telemetry.md and tasks/023-brain-top-production.md.
func renderSnapshot(w io.Writer, version string, procs ProcessReader, queues QueueReader, swaps SwapReader) {
	writeHeader(w, version)
	fmt.Fprintln(w)
	writeProcessPanel(w, procs)
	fmt.Fprintln(w)
	writeResourcePanel(w)
	fmt.Fprintln(w)
	writeQueuePanel(w, queues)
	fmt.Fprintln(w)
	writeSwapPanel(w, swaps)
}

// writeHeader renders the top header bar with version and timestamp.
func writeHeader(w io.Writer, version string) {
	ts := time.Now().Format("15:04:05")
	title := fmt.Sprintf("brain-top %s", version)
	separator := strings.Repeat("─", boxWidth-len(title)-len(ts)-4)
	fmt.Fprintf(w, "%s %s %s\n", bold(title), dim(separator), dim(ts))
}

// writeProcessPanel renders Panel 1: Agent process states.
func writeProcessPanel(w io.Writer, procs ProcessReader) {
	var lines []string

	// Column header.
	lines = append(lines, fmt.Sprintf("%-10s %-12s %-4s %-8s %-8s %-8s",
		bold("PID"), bold("STATE"), bold("PRI"), bold("MEM"), bold("VRAM"), bold("TASK")))

	rows := []process.Process{}
	if procs != nil {
		rows = procs.List()
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].PID < rows[j].PID })

	if len(rows) == 0 {
		lines = append(lines, dim("(no processes)"))
	} else {
		for _, p := range rows {
			stateStr := colorize(padRight(p.State.String(), 12), stateColor(p.State.String()))
			memStr := formatBytes(p.MemoryUsage)
			vramStr := formatBytes(p.VRAMUsage)
			lines = append(lines, fmt.Sprintf("%-10s %s %-4d %-8s %-8s %-8s",
				truncate(p.PID, 10), stateStr, p.Priority,
				memStr, vramStr, truncate(p.TaskID, 8)))
		}
	}

	writeBox(w, "Processes", boxWidth, lines)
}

// writeQueuePanel renders Panel 3: MLFQ queue depths.
func writeQueuePanel(w io.Writer, queues QueueReader) {
	var lines []string

	depths := [scheduler.NumLevels]int{}
	if queues != nil {
		depths = queues.Depths()
	}

	// Find max depth for scaling the bars.
	maxDepth := 1
	total := 0
	for _, d := range depths {
		if d > maxDepth {
			maxDepth = d
		}
		total += d
	}

	queueColors := [scheduler.NumLevels]string{ansiGreen, ansiCyan, ansiYellow, ansiRed}
	for i := 0; i < scheduler.NumLevels; i++ {
		label := fmt.Sprintf("Q%d", i)
		depthBar := ""
		barLen := 0
		if maxDepth > 0 {
			barLen = (depths[i] * 20) / maxDepth
		}
		for j := 0; j < barLen; j++ {
			depthBar += "█"
		}
		depthBar = colorize(depthBar, queueColors[i])
		lines = append(lines, fmt.Sprintf("  %s │ %-20s %d",
			bold(label), depthBar, depths[i]))
	}

	lines = append(lines, fmt.Sprintf("  %s   %d", dim("Total:"), total))

	writeBox(w, "MLFQ Queues", boxWidth, lines)
}

// writeResourcePanel renders Panel 2: GPU/RAM utilization bars.
func writeResourcePanel(w io.Writer) {
	var lines []string

	profile, err := hardware.ProbeAndClassify()
	if err != nil {
		lines = append(lines, fmt.Sprintf("  %s %v", colorize("probe failed:", ansiRed), err))
		writeBox(w, "Resources", boxWidth, lines)
		return
	}

	// Profile and backend info.
	lines = append(lines, fmt.Sprintf("  Profile: %s  Backend: %s  CPU: %s",
		bold(string(profile.Name)),
		colorize(string(profile.Probe.Backend), ansiCyan),
		dim(profile.Probe.CPUInfo)))

	lines = append(lines, "")

	// RAM bar — show total (free not available from ProbeResult, show total only).
	ramGiB := bytesToGiB(profile.Probe.TotalRAMBytes)
	lines = append(lines, fmt.Sprintf("  RAM:  %6.1f GiB total", ramGiB))

	// VRAM bar (if available).
	if profile.Probe.VRAMBytes > 0 {
		vramGiB := bytesToGiB(profile.Probe.VRAMBytes)
		lines = append(lines, fmt.Sprintf("  VRAM: %6.1f GiB total", vramGiB))
	} else {
		lines = append(lines, dim("  VRAM: (none detected)"))
	}

	writeBox(w, "Resources", boxWidth, lines)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

func bytesToGiB(b uint64) float64 {
	const gib = 1024 * 1024 * 1024
	return float64(b) / gib
}

// clearScreen writes an ANSI clear sequence to reset the terminal view.
// Works on all modern terminals: Windows Terminal (Win 10 1511+), Linux, macOS.
func clearScreen(w io.Writer) {
	fmt.Fprint(w, "\033[H\033[2J")
}
