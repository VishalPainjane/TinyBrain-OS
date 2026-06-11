package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/VishalPainjane/TinyBrain-OS/internal/hardware"
)

type probeJSON struct {
	Version string `json:"version"`
	Profile string `json:"profile"`
	RAMGiB  float64 `json:"ram_gib"`
	VRAMGiB float64 `json:"vram_gib"`
	Backend string `json:"backend"`
	CPUInfo string `json:"cpu_info"`
}

// runProbe prints hardware profile and returns an exit code.
func runProbe(w io.Writer, jsonOut bool) int {
	profile, err := hardware.ProbeAndClassify()
	if err != nil {
		fmt.Fprintf(os.Stderr, "probe failed: %v\n", err)
		return 1
	}

	if jsonOut {
		out := probeJSON{
			Version: Version,
			Profile: string(profile.Name),
			RAMGiB:  bytesToGiB(profile.Probe.TotalRAMBytes),
			VRAMGiB: bytesToGiB(profile.Probe.VRAMBytes),
			Backend: string(profile.Probe.Backend),
			CPUInfo: profile.Probe.CPUInfo,
		}
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(out); err != nil {
			fmt.Fprintf(os.Stderr, "encode json: %v\n", err)
			return 1
		}
		return 0
	}

	fmt.Fprintf(w, "TinyBrain %s | probe\n\n", Version)
	fmt.Fprintf(w, "  profile:  %s\n", profile.Name)
	fmt.Fprintf(w, "  ram:      %.1f GiB\n", bytesToGiB(profile.Probe.TotalRAMBytes))
	fmt.Fprintf(w, "  vram:     %.1f GiB\n", bytesToGiB(profile.Probe.VRAMBytes))
	fmt.Fprintf(w, "  backend:  %s\n", profile.Probe.Backend)
	fmt.Fprintf(w, "  cpu:      %s\n", profile.Probe.CPUInfo)
	return 0
}

func bytesToGiB(b uint64) float64 {
	const gib = 1024 * 1024 * 1024
	return float64(b) / gib
}
