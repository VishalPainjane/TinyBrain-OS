package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
)

// runAgentsList prints registered agents and returns an exit code.
func runAgentsList(w io.Writer) int {
	areg := registry.NewAgentRegistry()
	seedPath := filepath.Join("testdata", "fleet.yaml")
	
	if err := registry.LoadAgentsYAML(seedPath, areg); err != nil {
		fmt.Fprintf(os.Stderr, "open agent registry: %v\n", err)
		return 1
	}

	agentsList := areg.ListAgents()
	fmt.Fprintf(w, "TinyBrain %s | agents (%d)\n\n", Version, len(agentsList))
	if len(agentsList) == 0 {
		fmt.Fprintln(w, "  (no agents — verify testdata/fleet.yaml)")
		return 0
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tNAME\tMODEL")
	for _, a := range agentsList {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", a.ID, a.Name, a.ModelProfile)
	}
	_ = tw.Flush()
	return 0
}
