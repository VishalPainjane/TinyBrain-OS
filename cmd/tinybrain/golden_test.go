package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Golden files live in testdata/golden/. Regenerate with:
//
//	TB_UPDATE_GOLDEN=1 go test ./cmd/tinybrain/ -run TestGolden
func TestGolden_modelsListSeeded(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TB_MODELS_DB", filepath.Join(dir, "models.db"))
	t.Setenv("TB_MODELS_SEED", filepath.Join("..", "..", "testdata", "models.yaml"))

	var out bytes.Buffer
	if code := runModelsList(&out); code != 0 {
		t.Fatalf("runModelsList exit = %d", code)
	}

	assertGoldenText(t, "models_list_seeded.txt", out.Bytes())
}

func TestGolden_probeJSONNormalized(t *testing.T) {
	var out bytes.Buffer
	if code := runProbe(&out, true); code != 0 {
		t.Fatalf("runProbe exit = %d", code)
	}

	var payload probeJSON
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal probe json: %v", err)
	}

	normalized := probeJSON{
		Version: payload.Version,
		Profile: "<profile>",
		RAMGiB:  0,
		VRAMGiB: 0,
		Backend: "<backend>",
		CPUInfo: "<cpu>",
	}

	got, err := json.MarshalIndent(normalized, "", "  ")
	if err != nil {
		t.Fatalf("marshal normalized: %v", err)
	}
	got = append(got, '\n')

	assertGoldenText(t, "probe_normalized.json", got)
}

func assertGoldenText(t *testing.T, name string, got []byte) {
	t.Helper()

	goldenPath := filepath.Join("testdata", "golden", name)
	if os.Getenv("TB_UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
			t.Fatalf("mkdir golden: %v", err)
		}
		if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
		t.Logf("updated golden %s", goldenPath)
		return
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden %s: %v (run TB_UPDATE_GOLDEN=1 go test ./cmd/tinybrain/ -run TestGolden)", goldenPath, err)
	}

	if normalizeGoldenText(want) != normalizeGoldenText(got) {
		t.Fatalf("golden mismatch %s\n--- want ---\n%s--- got ---\n%s", name, want, got)
	}
}

func normalizeGoldenText(b []byte) string {
	return strings.ReplaceAll(string(b), "\r\n", "\n")
}

func TestGolden_doctorHeader(t *testing.T) {
	t.Setenv("TB_MODELS_DB", filepath.Join(t.TempDir(), "models.db"))

	var out bytes.Buffer
	if code := runDoctor(&out); code != 0 && code != 1 {
		t.Fatalf("runDoctor exit = %d", code)
	}

	firstLine, _, _ := strings.Cut(out.String(), "\n")
	want := "TinyBrain " + Version + " | doctor"
	if firstLine != want {
		t.Fatalf("doctor header = %q, want %q", firstLine, want)
	}
}
