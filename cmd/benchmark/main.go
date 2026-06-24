package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
)

func main() {
	mode := flag.String("mode", "monolith", "Benchmark mode: 'monolith' or 'swarm'")
	prompt := flag.String("prompt", "Write a comprehensive essay about the history of artificial intelligence.", "Test prompt")
	flag.Parse()

	if *mode != "monolith" && *mode != "swarm" {
		fmt.Fprintf(os.Stderr, "Invalid mode: %s\n", *mode)
		os.Exit(2)
	}

	if os.Getenv("CGO_ENABLED") == "0" {
		fmt.Fprintln(os.Stderr, "error: CGO is disabled — set CGO_ENABLED=1 and build llama.cpp for inference")
		os.Exit(1)
	}

	fmt.Printf("Starting benchmark: %s\n", *mode)
	code := runBenchmark(*mode, *prompt, os.Stdout, os.Stderr)
	os.Exit(code)
}

func openRegistry() (*registry.ModelRegistry, error) {
	dbPath := filepath.Join(".tinybrain", "models.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("create registry dir: %w", err)
	}
	return registry.NewBboltModelRegistry(dbPath, "")
}

func ngLayersFromEnv() (int32, bool) {
	raw := os.Getenv("TB_NGLAYERS")
	if raw == "" {
		return 0, false
	}
	v, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return 0, false
	}
	return int32(v), true
}
