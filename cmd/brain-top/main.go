// Command brain-top is a read-only terminal dashboard for TinyBrain OS.
// It displays four panels: process states, resource utilization, MLFQ queue
// depths, and swap monitor activity.
// See docs/architecture/telemetry.md and tasks/023-brain-top-production.md.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// Version tracks the brain-top release line.
const Version = "1.0.0"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	interval := time.Second
	once := true

	// Check NO_COLOR environment variable.
	if os.Getenv("NO_COLOR") != "" {
		colorEnabled = false
	}

	if len(args) > 0 {
		switch args[0] {
		case "snapshot":
			once = true
		case "watch":
			once = false
			if len(args) > 1 {
				d, err := time.ParseDuration(args[1])
				if err != nil {
					fmt.Fprintf(stderr, "invalid interval: %v\n", err)
					return 2
				}
				if d < 200*time.Millisecond {
					fmt.Fprintln(stderr, "interval must be >= 200ms")
					return 2
				}
				interval = d
			}
		case "version", "--version", "-v":
			fmt.Fprintf(stdout, "brain-top %s\n", Version)
			return 0
		case "help", "--help", "-h":
			printUsage(stdout)
			return 0
		default:
			fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
			printUsage(stderr)
			return 2
		}
	}

	procs := processTableReader{}
	queues := mlfqQueueReader{}
	var swaps SwapReader // nil — no swap events until wired to live kernel

	if once {
		renderSnapshot(stdout, Version, procs, queues, swaps)
		return 0
	}

	// Watch mode with graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	watchSignals(cancel)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Render immediately, then on each tick.
	clearScreen(stdout)
	renderSnapshot(stdout, Version, procs, queues, swaps)

	for {
		select {
		case <-ctx.Done():
			fmt.Fprintln(stdout, "\nbrain-top: shutting down")
			return 0
		case <-ticker.C:
			clearScreen(stdout)
			renderSnapshot(stdout, Version, procs, queues, swaps)
		}
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintf(w, `brain-top %s — TinyBrain OS process and resource dashboard

Usage:
  brain-top [snapshot]     One-shot dashboard render (default)
  brain-top watch [interval]  Live refresh (default 1s, minimum 200ms)
  brain-top version        Print version
  brain-top help           Print this help

Examples:
  brain-top                Snapshot of current state
  brain-top watch          Live dashboard, 1s refresh
  brain-top watch 500ms    Live dashboard, 500ms refresh

Environment:
  NO_COLOR=1    Disable ANSI colours

Press Ctrl+C to exit watch mode.

`, Version)
}
