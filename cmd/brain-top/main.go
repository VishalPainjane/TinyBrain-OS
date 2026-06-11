// Command brain-top is a read-only terminal dashboard prototype for TinyBrain OS.
package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Version tracks the brain-top prototype release line.
const Version = "0.1.0-proto"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	interval := time.Second
	once := true

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

	if once {
		renderSnapshot(stdout, Version, procs, queues)
		return 0
	}

	for {
		clearScreen(stdout)
		renderSnapshot(stdout, Version, procs, queues)
		time.Sleep(interval)
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintf(w, `brain-top %s — read-only process and queue dashboard (prototype)

Usage:
  brain-top [snapshot]
  brain-top watch [interval]

Examples:
  brain-top
  brain-top watch 1s

Notes:
  Prototype renders empty process/queue panels unless wired to a live kernel.
  Does not mutate scheduler or runtime state.

`, Version)
}
