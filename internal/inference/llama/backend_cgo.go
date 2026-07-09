//go:build cgo

package llama

import (
	"fmt"

	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
)

// cgoBackend routes all inference calls through static CGO bindings.
// On Linux/macOS this is the only backend.  On Windows it is used for the CPU
// path and for any platform where ggml-cuda.dll is not present.
//
// The actual native calls live in bindings_common.go (shared across CPU and
// CUDA static builds) and bindings_cpu.go / bindings_cuda.go (link flags).
type cgoBackend struct{}

func (b *cgoBackend) loadModel(path, modelID string, cfg LlamaConfig) error {
	return loadNativeModel(path, modelID, cfg)
}

func (b *cgoBackend) unloadModel(modelID string) error {
	return unloadNativeModel(modelID)
}

func (b *cgoBackend) generate(req runtime.GenerateRequest, cfg LlamaConfig) (string, uint32, generateStats, error) {
	return generateNativeTimed(req.ModelID, req.Prompt, cfg)
}

func (b *cgoBackend) getMetadata(modelID string) (runtime.ModelCapabilities, error) {
	// Not yet implemented for static CGO builds
	return runtime.ModelCapabilities{
		ModelID:            modelID,
		ChatTemplate:       "",
		SupportsMultimodal: false,
		SupportsTools:      false,
		SupportsGrammar:    false,
	}, nil
}

func (b *cgoBackend) formatChat(modelID string, messages []runtime.ChatMessage, opts runtime.FormatChatOpts) (string, string, error) {
	// Fallback implementation for static CGO builds
	fallbackTmpl := ""
	for _, m := range messages {
		fallbackTmpl += fmt.Sprintf("<|im_start|>%s\n%s<|im_end|>\n", m.Role, m.Content)
	}
	if opts.AddGenerationPrompt {
		fallbackTmpl += "<|im_start|>assistant\n"
	}
	return fallbackTmpl, "fallback_chatml_cgo", nil
}

func (b *cgoBackend) saveContext(modelID, ctxID string) error {
	return fmt.Errorf("%w: saveContext not implemented for static CGO backend", ErrNotImplemented)
}

func (b *cgoBackend) restoreContext(modelID, ctxID string) error {
	return fmt.Errorf("%w: restoreContext not implemented for static CGO backend", ErrNotImplemented)
}
