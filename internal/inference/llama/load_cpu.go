//go:build cgo && !cuda && !rocm && !metal && !vulkan

// Package llama — CPU static CGO load/unload shim.
//
// These functions exist only for the integration test path
// (generate_integration_test.go / load_integration_test.go) which calls
// loadNativeModel / unloadNativeModel directly.
// The provider's public API routes through backend_cgo.go via selectBackend().
package llama
