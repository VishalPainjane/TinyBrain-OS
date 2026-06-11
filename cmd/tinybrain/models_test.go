package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunModelsList_empty(t *testing.T) {
	t.Setenv("TB_MODELS_DB", filepath.Join(t.TempDir(), "models.db"))
	t.Setenv("TB_MODELS_SEED", "")

	var out bytes.Buffer
	code := runModelsList(&out)
	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "models (0)") {
		t.Fatalf("output = %q", out.String())
	}
}

func TestRunModelsList_seeded(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TB_MODELS_DB", filepath.Join(dir, "models.db"))
	t.Setenv("TB_MODELS_SEED", filepath.Join("..", "..", "testdata", "models.yaml"))

	var out bytes.Buffer
	code := runModelsList(&out)
	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "tinyllama-q4") {
		t.Fatalf("output = %q, want seeded model", out.String())
	}
}
