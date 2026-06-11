package registry

import (
	"os"
	"path/filepath"
	"testing"
)

func FuzzLoadModelsYAML(f *testing.F) {
	f.Add([]byte("models:\n  - id: m1\n    path: /p/m1.gguf\n"))
	f.Add([]byte("not yaml"))
	f.Add([]byte("models:\n  - path: /no-id.gguf\n"))
	f.Add([]byte(""))

	f.Fuzz(func(t *testing.T, data []byte) {
		dir := t.TempDir()
		path := filepath.Join(dir, "models.yaml")
		if err := os.WriteFile(path, data, 0o644); err != nil {
			t.Fatalf("write temp yaml: %v", err)
		}
		store := NewModelRegistry()
		_ = LoadModelsYAML(path, store)
	})
}
