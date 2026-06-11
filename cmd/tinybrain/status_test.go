package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunStatus(t *testing.T) {
	t.Setenv("TB_MODELS_DB", filepath.Join(t.TempDir(), "models.db"))

	var out bytes.Buffer
	code := runStatus(&out)
	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	body := out.String()
	for _, want := range []string{"status", "profile:", "registry:"} {
		if !strings.Contains(body, want) {
			t.Fatalf("output missing %q: %s", want, body)
		}
	}
}
