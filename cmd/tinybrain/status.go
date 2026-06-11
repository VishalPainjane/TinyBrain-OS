package main

import (
	"fmt"
	"io"
	"os"

	"github.com/VishalPainjane/TinyBrain-OS/internal/hardware"
)

// runStatus prints a read-only runtime snapshot and returns an exit code.
func runStatus(w io.Writer) int {
	fmt.Fprintf(w, "TinyBrain %s | status\n\n", Version)

	profile, err := hardware.ProbeAndClassify()
	if err != nil {
		fmt.Fprintf(w, "  profile:  (probe failed: %v)\n", err)
	} else {
		fmt.Fprintf(w, "  profile:  %s\n", profile.Name)
		fmt.Fprintf(w, "  backend:  %s\n", profile.Probe.Backend)
		fmt.Fprintf(w, "  ram:      %.1f GiB\n", bytesToGiB(profile.Probe.TotalRAMBytes))
		fmt.Fprintf(w, "  vram:     %.1f GiB\n", bytesToGiB(profile.Probe.VRAMBytes))
	}

	if cgoEnabled() {
		fmt.Fprintln(w, "  cgo:      enabled")
	} else {
		fmt.Fprintln(w, "  cgo:      disabled")
	}

	fmt.Fprintf(w, "  llama:    %s\n", llamaLibDir())
	fmt.Fprintf(w, "  registry: %s\n", modelsDBPath())

	reg, err := openRegistry()
	if err != nil {
		fmt.Fprintf(os.Stderr, "open registry: %v\n", err)
		return 1
	}
	defer reg.Close()
	fmt.Fprintf(w, "  models:   %d registered\n", len(reg.ListModels()))
	return 0
}
