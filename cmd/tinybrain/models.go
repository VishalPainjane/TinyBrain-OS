package main

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// runModelsList prints registered models and returns an exit code.
func runModelsList(w io.Writer) int {
	reg, err := openRegistry()
	if err != nil {
		fmt.Fprintf(os.Stderr, "open registry: %v\n", err)
		return 1
	}
	defer reg.Close()

	models := reg.ListModels()
	fmt.Fprintf(w, "TinyBrain %s | models (%d)\n\n", Version, len(models))
	if len(models) == 0 {
		fmt.Fprintln(w, "  (no models — set TB_MODELS_SEED or register via API)")
		return 0
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tQUANT\tPATH")
	for _, m := range models {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", m.ID, m.Quantization, m.Path)
	}
	_ = tw.Flush()
	return 0
}
