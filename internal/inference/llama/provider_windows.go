//go:build windows && cgo && !cuda

// Package llama — Windows provider backend selector.
//
// On Windows the provider probes for ggml-cuda.dll at the DLL search path.
// If found, it uses the dynamic DLL backend (no CGO required for CUDA path).
// If not found, it falls back to the static CGO CPU backend.
//
// This file replaces the compile-time -tags cuda split on Windows with a
// single binary that decides the backend at runtime.  See ADR-006.
package llama

import (
	"os"
	"path/filepath"
)

// selectBackend returns the best available backend for this process.
// Priority:
//  1. Windows GPU (dynamic DLL) if ggml-cuda.dll exists alongside the binary.
//  2. CGO CPU (static link) — always available when built with CGO_ENABLED=1.
func selectBackend() backend {
	dllDir := dllSearchDir()
	dynBackend, err := NewWindowsDynamicBackend(dllDir)
	if err == nil {
		return dynBackend
	}
	// DLL not present or failed to load — fall back to CPU CGO backend.
	// The error is intentionally silenced here; the caller can check
	// hardware.Probe() to know whether a GPU was expected.
	return &cgoBackend{}
}

// dllSearchDir returns the directory to probe for llama.dll / ggml-cuda.dll.
// Order of precedence:
//  1. LLAMACPP_DLL_DIR environment variable (explicit override)
//  2. Directory alongside the running executable
func dllSearchDir() string {
	if override := os.Getenv("LLAMACPP_DLL_DIR"); override != "" {
		return override
	}
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}
