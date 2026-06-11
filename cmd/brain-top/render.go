package main

import (
	"fmt"
	"io"
	"sort"
	"strings"

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

// renderSnapshot writes a read-only brain-top dashboard to w.
// See docs/architecture/telemetry.md and tasks/013-brain-top.md.
func renderSnapshot(w io.Writer, version string, procs ProcessReader, queues QueueReader) {
	fmt.Fprintf(w, "brain-top %s | snapshot\n\n", version)

	writeProcessPanel(w, procs)
	fmt.Fprintln(w)
	writeQueuePanel(w, queues)
	fmt.Fprintln(w)
	writeResourcePanel(w)
}

func writeProcessPanel(w io.Writer, procs ProcessReader) {
	fmt.Fprintln(w, "[Processes]")
	fmt.Fprintf(w, "%-10s %-12s %-4s %-10s %-10s\n", "PID", "STATE", "PRI", "KV", "TASK")

	rows := []process.Process{}
	if procs != nil {
		rows = procs.List()
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].PID < rows[j].PID })

	if len(rows) == 0 {
		fmt.Fprintln(w, "(none)")
		return
	}
	for _, p := range rows {
		fmt.Fprintf(w, "%-10s %-12s %-4d %-10s %-10s\n",
			p.PID, p.State.String(), p.Priority, truncate(p.KVCacheID, 10), truncate(p.TaskID, 10))
	}
}

func writeQueuePanel(w io.Writer, queues QueueReader) {
	fmt.Fprintln(w, "[Queues]")
	depths := [scheduler.NumLevels]int{}
	if queues != nil {
		depths = queues.Depths()
	}
	fmt.Fprintf(w, "Q0 Q1 Q2 Q3\n")
	fmt.Fprintf(w, "%2d %2d %2d %2d\n",
		depths[0], depths[1], depths[2], depths[3])
}

func writeResourcePanel(w io.Writer) {
	fmt.Fprintln(w, "[Resources]")
	profile, err := hardware.ProbeAndClassify()
	if err != nil {
		fmt.Fprintf(w, "probe failed: %v\n", err)
		return
	}
	fmt.Fprintf(w, "profile: %-10s backend: %s\n", profile.Name, profile.Probe.Backend)
	fmt.Fprintf(w, "ram:     %.1f GiB   vram: %.1f GiB\n",
		bytesToGiB(profile.Probe.TotalRAMBytes), bytesToGiB(profile.Probe.VRAMBytes))
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

// clearScreen writes an ANSI clear sequence when w is a TTY-friendly writer.
func clearScreen(w io.Writer) {
	fmt.Fprint(w, strings.Repeat("\n", 2))
	fmt.Fprint(w, "\033[H\033[2J")
}
