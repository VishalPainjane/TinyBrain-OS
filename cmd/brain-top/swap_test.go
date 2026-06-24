package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

type mockSwapReader struct {
	events []SwapRecord
}

func (m mockSwapReader) SwapEvents() []SwapRecord {
	return m.events
}

func TestWriteSwapPanel_NilReader(t *testing.T) {
	t.Parallel()

	colorEnabled = false
	defer func() { colorEnabled = true }()

	var buf bytes.Buffer
	writeSwapPanel(&buf, nil)
	out := buf.String()

	if !strings.Contains(out, "Swap Monitor") {
		t.Error("missing panel title")
	}
	if !strings.Contains(out, "(no swap activity)") {
		t.Errorf("expected no-activity message, got:\n%s", out)
	}
}

func TestWriteSwapPanel_EmptyReader(t *testing.T) {
	t.Parallel()

	colorEnabled = false
	defer func() { colorEnabled = true }()

	var buf bytes.Buffer
	writeSwapPanel(&buf, mockSwapReader{})
	if !strings.Contains(buf.String(), "(no swap activity)") {
		t.Fatalf("expected no-activity message:\n%s", buf.String())
	}
}

func TestWriteSwapPanel_WithEvents(t *testing.T) {
	t.Parallel()

	colorEnabled = false
	defer func() { colorEnabled = true }()

	events := []SwapRecord{
		{PID: "pid-1", From: "VRAM", To: "RAM", Timestamp: time.Now(), SizeBytes: 1024 * 1024 * 512},
		{PID: "pid-2", From: "RAM", To: "SSD", Timestamp: time.Now(), SizeBytes: 1024 * 1024 * 1024 * 2},
	}

	var buf bytes.Buffer
	writeSwapPanel(&buf, mockSwapReader{events: events})
	out := buf.String()

	for _, want := range []string{"pid-1", "pid-2", "VRAM", "RAM", "SSD", "512 MiB", "2.0 GiB"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n%s", want, out)
		}
	}
}

func TestFormatBytes(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input uint64
		want  string
	}{
		{0, "0 B"},
		{1023, "1023 B"},
		{1024 * 1024, "1 MiB"},
		{1024 * 1024 * 1024, "1.0 GiB"},
		{1024 * 1024 * 1024 * 2, "2.0 GiB"},
		{1024 * 1024 * 512, "512 MiB"},
	}

	for _, tc := range cases {
		got := formatBytes(tc.input)
		if got != tc.want {
			t.Errorf("formatBytes(%d) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
