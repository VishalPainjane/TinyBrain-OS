package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestRunProbe_table(t *testing.T) {
	var out bytes.Buffer
	code := runProbe(&out, false)
	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	if !strings.Contains(out.String(), "profile:") {
		t.Fatalf("output = %q, want profile line", out.String())
	}
}

func TestRunProbe_json(t *testing.T) {
	var out bytes.Buffer
	code := runProbe(&out, true)
	if code != 0 {
		t.Fatalf("exit = %d, want 0", code)
	}
	var payload probeJSON
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if payload.Version != Version {
		t.Fatalf("version = %q, want %s", payload.Version, Version)
	}
	if payload.Profile == "" {
		t.Fatal("profile is empty")
	}
}
