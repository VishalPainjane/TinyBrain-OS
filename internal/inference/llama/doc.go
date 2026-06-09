// Package llama implements a local GGUF inference adapter via llama.cpp.
// CGO and llama.cpp imports are confined to this package (INV-008).
// CPU is the default backend; optional CUDA offload requires -tags cuda and
// a llama.cpp build with GGML_CUDA=ON (see planning/decisions/009d-architecture-review.md).
// See docs/architecture/inference-lifecycle.md and planning/decisions/009b-architecture-review.md.
package llama
