//go:build cgo && cuda && !rocm && !metal && !vulkan

// Package llama — CUDA static CGO load/unload shim (Linux CUDA builds).
//
// On Windows the CUDA path uses the dynamic DLL backend (backend_windows_dynamic.go).
// On Linux, nvcc uses GCC as host compiler so static CGO is used directly.
// The provider's public API routes through backend_cgo.go via selectBackend().
package llama
