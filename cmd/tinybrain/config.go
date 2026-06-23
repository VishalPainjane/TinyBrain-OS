package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
)

// Version is the CLI release version (matches latest shipped tag).
const Version = "0.8.0"

// dataDir returns the TinyBrain config directory (~/.tinybrain).
func dataDir() string {
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".tinybrain")
	}
	return ".tinybrain"
}

// modelsDBPath returns the model registry database path.
func modelsDBPath() string {
	if p := os.Getenv("TB_MODELS_DB"); p != "" {
		return p
	}
	return filepath.Join(dataDir(), "models.db")
}

// modelsSeedPath returns the optional YAML seed path for an empty registry.
func modelsSeedPath() string {
	return os.Getenv("TB_MODELS_SEED")
}

// llamaLibDir returns the directory containing llama.cpp shared libraries.
func llamaLibDir() string {
	if d := os.Getenv("TB_LLAMA_LIB_DIR"); d != "" {
		return d
	}
	return filepath.Join("third_party", "llama.cpp", "build", "bin")
}

// ngLayersFromEnv reads TB_NGLAYERS when set; ok is false when unset.
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

// openRegistry opens or creates the bbolt model registry.
func openRegistry() (*registry.ModelRegistry, error) {
	dbPath := modelsDBPath()
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("create registry dir: %w", err)
	}
	return registry.NewBboltModelRegistry(dbPath, modelsSeedPath())
}

// cgoEnabled reports whether CGO is enabled for this process.
func cgoEnabled() bool {
	return os.Getenv("CGO_ENABLED") != "0"
}
