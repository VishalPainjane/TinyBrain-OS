//go:build !windows && cgo

// Package llama — non-Windows provider backend selector.
//
// On Linux and macOS, llama.cpp is always linked via static CGO (gcc/clang).
// nvcc uses GCC as the host compiler on Linux, so the static CGO path works
// for CUDA builds on Linux too.  There is no dynamic loading needed.
package llama

// selectBackend always returns the CGO backend on non-Windows platforms.
func selectBackend() backend {
	return &cgoBackend{}
}
