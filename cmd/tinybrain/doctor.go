package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

// runDoctor prints startup diagnostics and returns an exit code.
func runDoctor(w io.Writer) int {
	fmt.Fprintf(w, "TinyBrain %s | doctor\n\n", Version)

	var results []checkResult

	goVer := runtime.Version()
	if goVer != "" {
		results = append(results, printCheck(w, checkOK, "Go %s (%s/%s)", goVer, runtime.GOOS, runtime.GOARCH))
	}

	if cgoEnabled() {
		results = append(results, printCheck(w, checkOK, "CGO enabled"))
	} else {
		results = append(results, printCheck(w, checkWarn, "CGO disabled — inference run requires CGO_ENABLED=1"))
	}

	libDir := llamaLibDir()
	if llamaBuildPresent(libDir) {
		results = append(results, printCheck(w, checkOK, "llama.cpp build: %s", libDir))
	} else {
		results = append(results, printCheck(w, checkWarn, "llama.cpp libraries not found in %s", libDir))
	}

	dbPath := modelsDBPath()
	reg, err := openRegistry()
	if err != nil {
		results = append(results, printCheck(w, checkFail, "model registry %s: %v", dbPath, err))
	} else {
		defer reg.Close()
		models := reg.ListModels()
		results = append(results, printCheck(w, checkOK, "model registry: %s (%d models)", dbPath, len(models)))
		if len(models) == 0 {
			results = append(results, printCheck(w, checkWarn, "no models registered — set TB_MODELS_SEED or register models"))
		} else {
			for _, m := range models {
				if _, err := os.Stat(m.Path); err != nil {
					results = append(results, printCheck(w, checkWarn, "model %q path missing: %s", m.ID, m.Path))
				}
			}
		}
	}

	if gguf := os.Getenv("TB_TEST_GGUF_PATH"); gguf != "" {
		if _, err := os.Stat(gguf); err != nil {
			results = append(results, printCheck(w, checkWarn, "TB_TEST_GGUF_PATH not accessible: %s", gguf))
		} else {
			results = append(results, printCheck(w, checkOK, "TB_TEST_GGUF_PATH set"))
		}
	}

	return exitCodeForChecks(results)
}

func llamaBuildPresent(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.Contains(strings.ToLower(e.Name()), "llama") {
			return true
		}
	}
	return false
}
