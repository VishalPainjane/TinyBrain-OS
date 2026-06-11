package main

import (
	"os"
	"testing"
)

func FuzzNgLayersFromEnv(f *testing.F) {
	seeds := []string{"", "-1", "0", "32", "not-a-number", "999999999999999999999"}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, raw string) {
		if err := os.Setenv("TB_NGLAYERS", raw); err != nil {
			return
		}
		t.Cleanup(func() { _ = os.Unsetenv("TB_NGLAYERS") })
		_, _ = ngLayersFromEnv()
	})
}

func FuzzModelsDBPath(f *testing.F) {
	f.Add("")
	f.Add("/tmp/custom/models.db")
	f.Add("relative/path.db")

	f.Fuzz(func(t *testing.T, path string) {
		if err := os.Setenv("TB_MODELS_DB", path); err != nil {
			return
		}
		t.Cleanup(func() { _ = os.Unsetenv("TB_MODELS_DB") })
		_ = modelsDBPath()
	})
}
