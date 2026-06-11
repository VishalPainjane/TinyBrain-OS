package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunDoctor(t *testing.T) {
	t.Setenv("TB_MODELS_DB", t.TempDir()+"/models.db")
	var out bytes.Buffer
	code := runDoctor(&out)
	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	body := out.String()
	if !strings.Contains(body, "doctor") {
		t.Fatalf("output missing doctor header: %q", body)
	}
	if !strings.Contains(body, "[ok]") && !strings.Contains(body, "[warn]") {
		t.Fatalf("output missing check lines: %q", body)
	}
}
