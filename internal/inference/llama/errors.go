package llama

import "errors"

var (
	// ErrCGODisabled indicates CGO is off or llama.cpp is not linked.
	ErrCGODisabled = errors.New("llama.cpp CGO backend disabled")
	// ErrNotImplemented indicates a port method is not implemented in this task.
	ErrNotImplemented = errors.New("not implemented")
	// ErrPathInaccessible indicates the GGUF path cannot be accessed.
	ErrPathInaccessible = errors.New("model path not accessible")
)
